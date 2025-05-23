FROM golang:1.24 as builder

WORKDIR /build
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o app-bin .

FROM alpine:3.21.3

WORKDIR /app
COPY --from=builder /build/app-bin /app/app-bin

CMD ["/app/app-bin"]
