# Step 1: Build stage
FROM golang:1.21.6-alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chatops ./cmd

# Step 2: Runtime stage
FROM alpine:latest

ENV TZ=Asia/Jakarta
RUN apk --no-cache add tzdata 
RUN apk --no-cache add ca-certificates 
RUN apk --no-cache add postgresql-client

ARG USER_ID=10001
ARG GROUP_ID=10001
RUN addgroup -g $GROUP_ID app && adduser -D -u $USER_ID -G app app

COPY --from=builder /app/chatops /home/app/chatops

RUN chown -R app:app /home/app/chatops

USER app

WORKDIR /home/app

ENTRYPOINT ["./chatops"]