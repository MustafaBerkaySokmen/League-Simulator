# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/league && go build -o /app/league

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/league ./league
COPY web/ ./web/
EXPOSE 8080
CMD ["./league"]
