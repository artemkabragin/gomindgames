FROM golang:1.25.2-alpine AS builder
WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod tidy

COPY . .

RUN ls -la /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main ./cmd/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8081

# Command to run
CMD ["./main"]