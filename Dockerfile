FROM golang:1.20-alpine as builder
RUN apk add --no-cache build-base
WORKDIR /build
COPY . .

RUN CGO_ENABLED=1 go build bonefabric/adviser/cmd/adviser

FROM alpine
WORKDIR /app
COPY --from=builder /build/adviser /app
COPY --from=builder /build/config.yaml /app

CMD ["./adviser"]
