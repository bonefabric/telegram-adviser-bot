FROM golang:1.20-alpine as builder
RUN apk add --no-cache build-base
WORKDIR /build
COPY . .

RUN CGO_ENABLED=1 go build -o app .


FROM alpine
WORKDIR /app
COPY --from=builder /build/app /app
COPY --from=builder /build/config.yaml /app

CMD ["./app"]
