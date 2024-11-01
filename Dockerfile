FROM golang:1.19.7 AS  build-env
ARG DIR
ENV GOPROXY=https://goproxy.cn
WORKDIR /data
COPY . /data/
RUN cd /data/infra/services/rubick && go mod tidy
RUN cd /data/infra/services/rubick/${DIR} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM registry.cn-hangzhou.aliyuncs.com/meetwhale/wop-debian-ca-certificates:3.0
ARG DIR
WORKDIR /data
COPY  --from=build-env /data/infra/services/rubick/${DIR}/main /data
EXPOSE 8981
CMD /data/main