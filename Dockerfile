FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pixoo64 .


FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Europe/Moscow

WORKDIR /app

COPY --from=builder /app/pixoo64 .
COPY --from=builder /app/cache ./cache
COPY --from=builder /app/static ./static
COPY --from=builder /app/mocks ./mocks

CMD ["./pixoo64"]
