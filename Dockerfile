FROM golang:1.16-alpine AS builder
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV GOPRIVATE=""
ENV GOPROXY="https://goproxy.cn,direct"
ENV GOSUMDB="sum.golang.google.cn"
WORKDIR /root/edgex-thingsboard/

# The main mirrors are giving us timeout issues on builds periodically.
# So we can try these.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.aliyun.com/g' -i~ /etc/apk/repositories
RUN apk update && apk add zeromq-dev libsodium-dev pkgconfig build-base git

ADD . .
RUN go mod download \
    && go test --cover $(go list ./... | grep -v /vendor/) \
    && go build -o main cmd/main.go

FROM alpine
WORKDIR /root/
ENV TZ Asia/Shanghai

# The main mirrors are giving us timeout issues on builds periodically.
# So we can try these.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.aliyun.com/g' -i~ /etc/apk/repositories
RUN apk --no-cache add zeromq

COPY --from=builder /root/edgex-thingsboard/main edgex-thingsboard
COPY --from=builder /root/edgex-thingsboard/cmd/res/docker/configuration.toml res/configuration.toml
RUN chmod +x edgex-thingsboard

ENTRYPOINT ["/root/edgex-thingsboard"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/root/res"]
