# Step 1: Build stage
FROM golang:1.21.6-alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chatops ./cmd

# Step 2: Runtime stage
FROM alpine:latest

WORKDIR /home/app
COPY --from=builder /app/chatops /home/app/chatops

ENTRYPOINT ["./chatops"]