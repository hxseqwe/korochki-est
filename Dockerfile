FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/server/main.go

FROM node:18-alpine AS frontend-builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM alpine:latest

RUN apk add --no-cache postgresql15-client

WORKDIR /root/

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/build ./frontend/build
COPY migrations ./migrations

EXPOSE 8080

CMD ["./main"]