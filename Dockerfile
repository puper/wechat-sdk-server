FROM golang as builder
WORKDIR     /code
COPY        . .
RUN mkdir -p /app/data/logs /app/config
COPY config /app/config

RUN CGO_ENABLED=0 go build -mod=vendor -o /app/wechat-sdk-server -v .

FROM alpine:latest
COPY --from=builder /app /app
ENTRYPOINT [ "/app/wechat-sdk-server", "--config", "/app/config/config.toml", "serve"]