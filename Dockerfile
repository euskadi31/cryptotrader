FROM golang:latest as builder
WORKDIR /go/src/github.com/euskadi31/cryptotrader/
COPY . .
RUN CGO_ENABLED=0 go build .

FROM alpine:latest
LABEL maintainer "axel@etcheverry.biz"
ENV PORT 8080
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /go/src/github.com/euskadi31/cryptotrader/cryptotrader .
COPY --from=builder /go/src/github.com/euskadi31/cryptotrader/config.yml.dist /etc/cryptotrader/config.yml
HEALTHCHECK --interval=1m --timeout=3s CMD curl -f http://localhost:${PORT}/health > /dev/null 2>&1 || exit 1
EXPOSE ${PORT}
VOLUME /var/lib/cryptotrader
CMD [ "./cryptotrader" ]
