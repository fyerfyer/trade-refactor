FROM golang:1.23.1 AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main/main.go

FROM alpine:latest
COPY --from=builder /usr/src/app/payment ./payment

RUN apk add --no-cache curl iputils
CMD ["./payment"]
