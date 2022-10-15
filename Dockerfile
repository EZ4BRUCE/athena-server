FROM golang:1.17.13-alpine

LABEL MAINTAINER="EZ4BRUCE@lhy122786302@gmail.com"

# 为镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/build
WORKDIR /go/src/athena-server
# 将代码复制到容器中
COPY . .
RUN go mod tidy -compat=1.17

# go generate 编译前自动执行代码
# go env 查看go的环境变量
# go build -o athena-server . 打包项目生成文件名为athena-server的二进制文件
RUN go generate && go env && go build -o athena-server .


FROM alpine:latest

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

LABEL MAINTAINER="EZ4BRUCE@lhy122786302@gmail.com"
WORKDIR /go/src/athena-server

# 把/go/src/gin-vue-admin整个文件夹的文件到当前工作目录
COPY --from=0 /go/src/athena-server ./

EXPOSE 8880

ENTRYPOINT ./athena-server 