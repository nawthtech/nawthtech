FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY  ..

RUN go build -ldflags="-w -s" -o nawthtech-backend .

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/nawthtech-backend ./

ENTRYPOINT ["/app/nawthtech-backend"]