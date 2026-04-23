FROM golang:1.25.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod /app/
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bayarwoy .

FROM alpine:3.21

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bayarwoy /app/bayarwoy

ENV PORT=8080

EXPOSE 8080

USER appuser

ENTRYPOINT ["/app/bayarwoy"]
