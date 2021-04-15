FROM golang:1.16-alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV GOPRIVATE=""
ENV GOPROXY="https://goproxy.cn,direct"
ENV GOSUMDB="sum.golang.google.cn"
WORKDIR /root/edgex-control-agent/

ADD . .
RUN go mod download \
    && go test --cover $(go list ./... | grep -v /vendor/) \
    && go build -o main cmd/agent/main.go

FROM alpine
WORKDIR /root/
ENV TZ Asia/Shanghai

COPY --from=builder /root/edgex-control-agent/main edgex-control-agent
COPY --from=builder /root/edgex-control-agent/cmd/agent/res/docker/configuration.toml res/configuration.toml
RUN chmod +x edgex-control-agent

ENTRYPOINT ["/root/edgex-control-agent"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/root/res"]
