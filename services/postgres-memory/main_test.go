package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/openai/openai-go"
)

func setupTestServer(t *testing.T) *Server {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database tests")
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	server, err := New(dsn, "messages", logger)
	if err != nil {
		t.Fatalf("create test server: %v", err)
	}

	// Clean test data
	_, err = server.db.Exec("DELETE FROM messages")
	if err != nil {
		t.Fatalf("clean test database: %v", err)
	}

	return server
}

func TestServer_AddMessage(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	ctx := context.Background()
	sessionID := "test-session"
	message := Message(openai.UserMessage("Hello, world!"))

	err := server.addMessage(ctx, sessionID, message)
	if err != nil {
		t.Fatalf("addMessage failed: %v", err)
	}

	var count int
	err = server.db.QueryRow("SELECT COUNT(*) FROM messages WHERE session_id = $1", sessionID).Scan(&count)
	if err != nil {
		t.Fatalf("count messages: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 message, got %d", count)
	}
}

func TestServer_GetMessages(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	ctx := context.Background()
	sessionID := "test-session"

	messages := []Message{
		Message(openai.UserMessage("First message")),
		Message(openai.AssistantMessage("Second message")),
	}

	for _, msg := range messages {
		if err := server.addMessage(ctx, sessionID, msg); err != nil {
			t.Fatalf("addMessage failed: %v", err)
		}
	}

	retrieved, err := server.getMessages(ctx, sessionID)
	if err != nil {
		t.Fatalf("getMessages failed: %v", err)
	}

	if len(retrieved) != len(messages) {
		t.Errorf("expected %d messages, got %d", len(messages), len(retrieved))
	}
}

func TestServer_HTTPHandlers(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	router := server.routes().(*mux.Router)

	t.Run("AddMessage", func(t *testing.T) {
		req := struct {
			Message Message `json:"message"`
		}{Message: Message(openai.UserMessage("Test message"))}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPut, "/message/test-session", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq = mux.SetURLVars(httpReq, map[string]string{"uid": "test-session"})

		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("GetMessages", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/message/test-session", nil)
		httpReq = mux.SetURLVars(httpReq, map[string]string{"uid": "test-session"})

		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response struct {
			Messages []Message `json:"messages"`
		}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Messages) == 0 {
			t.Error("expected at least one message")
		}
	})

	t.Run("Health", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if w.Body.String() != "OK" {
			t.Errorf("expected 'OK', got %q", w.Body.String())
		}
	})
}