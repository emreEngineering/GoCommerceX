FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

ARG SERVICE_PATH
ARG SERVICE_NAME

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./$SERVICE_PATH/cmd/$SERVICE_NAME

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -g '' appuser

WORKDIR /app

RUN mkdir -p /src

COPY --from=builder /out/service /app/service
COPY --from=builder /src/auth-service/migrations /src/auth-service/migrations
COPY --from=builder /src/user-service/migrations /src/user-service/migrations
COPY --from=builder /src/product-service/migrations /src/product-service/migrations
COPY --from=builder /src/inventory-service/migrations /src/inventory-service/migrations
COPY --from=builder /src/order-service/migrations /src/order-service/migrations
COPY --from=builder /src/payment-service/migrations /src/payment-service/migrations
COPY --from=builder /src/notification-service/migrations /src/notification-service/migrations

USER appuser

ENTRYPOINT ["/app/service"]
