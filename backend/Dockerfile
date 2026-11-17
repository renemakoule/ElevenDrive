# Étape 1 : Compilation du binaire Go
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache des dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copie du code source backend
COPY . .

# Compilation statique ultra-légère
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Étape 2 : Image d'exécution minimale
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
