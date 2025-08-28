# PostgreSQL Memory Service API Tests
# Test the complete PostgreSQL Memory service functionality

# Test health endpoint
GET http://postgres-memory/health
HTTP 200
[Asserts]
body == "OK"

# Test adding a single message
PUT http://postgres-memory/message/test-session
Content-Type: application/json
{
  "message": {
    "role": "user",
    "content": "Hello, this is a test message!"
  }
}
HTTP 200

# Test retrieving messages from session
GET http://postgres-memory/message/test-session
HTTP 200
[Asserts]
jsonpath "$.messages" count == 1
jsonpath "$.messages[0].role" == "user"
jsonpath "$.messages[0].content" == "Hello, this is a test message!"

# Test adding multiple messages in bulk
PUT http://postgres-memory/messages/bulk-session
Content-Type: application/json
{
  "messages": [
    {
      "role": "user",
      "content": "First message in bulk"
    },
    {
      "role": "assistant", 
      "content": "Second message in bulk"
    },
    {
      "role": "user",
      "content": "Third message in bulk"
    }
  ]
}
HTTP 200

# Test retrieving bulk messages
GET http://postgres-memory/message/bulk-session
HTTP 200
[Asserts]
jsonpath "$.messages" count == 3
jsonpath "$.messages[0].role" == "user"
jsonpath "$.messages[0].content" == "First message in bulk"
jsonpath "$.messages[1].role" == "assistant"
jsonpath "$.messages[1].content" == "Second message in bulk"
jsonpath "$.messages[2].role" == "user"
jsonpath "$.messages[2].content" == "Third message in bulk"

# Test session isolation - different session should have no messages
GET http://postgres-memory/message/empty-session
HTTP 200
[Asserts]
jsonpath "$.messages" == null

# Test adding another message to original session
PUT http://postgres-memory/message/test-session
Content-Type: application/json
{
  "message": {
    "role": "assistant",
    "content": "Response to the test message"
  }
}
HTTP 200

# Verify session now has two messages in correct order
GET http://postgres-memory/message/test-session
HTTP 200
[Asserts]
jsonpath "$.messages" count == 2
jsonpath "$.messages[0].role" == "user"
jsonpath "$.messages[0].content" == "Hello, this is a test message!"
jsonpath "$.messages[1].role" == "assistant"
jsonpath "$.messages[1].content" == "Response to the test message"

# Test with complex message structure (tool calls, etc.)
PUT http://postgres-memory/message/complex-session
Content-Type: application/json
{
  "message": {
    "role": "assistant",
    "content": null,
    "tool_calls": [
      {
        "id": "call_123",
        "type": "function",
        "function": {
          "name": "get_weather",
          "arguments": "{\"location\": \"San Francisco\"}"
        }
      }
    ]
  }
}
HTTP 200

# Verify complex message structure is preserved
GET http://postgres-memory/message/complex-session
HTTP 200
[Asserts]
jsonpath "$.messages" count == 1
jsonpath "$.messages[0].role" == "assistant"
jsonpath "$.messages[0].content" == null
jsonpath "$.messages[0].tool_calls" exists
jsonpath "$.messages[0].tool_calls[0].id" == "call_123"
jsonpath "$.messages[0].tool_calls[0].function.name" == "get_weather"

# Test complete - all core functionality validated