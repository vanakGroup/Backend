FROM golang:1.26-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/hub.megan.ir/g' /etc/apk/repositories || true
RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://hub.megan.ir/go,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vanak-backend ./cmd/app/main.go

FROM alpine:3.19

RUN sed -i 's/dl-cdn.alpinelinux.org/hub.megan.ir/g' /etc/apk/repositories || true
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/vanak-backend .
COPY --from=builder /app/.env .

COPY internal/infrastructure/jwt/ ./internal/infrastructure/jwt/

EXPOSE 8080

CMD ["./vanak-backend"]