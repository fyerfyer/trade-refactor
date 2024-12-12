FROM golang:1.23.1 AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main/main.go

FROM scratch
COPY --from=builder /usr/src/app/order ./order
CMD ["./order"]
