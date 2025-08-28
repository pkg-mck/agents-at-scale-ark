FROM golang:1.24-bullseye AS builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 go build -a -o postgres-memory ./main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/postgres-memory  .
USER 65532:65532

ENTRYPOINT ["/postgres-memory"]
