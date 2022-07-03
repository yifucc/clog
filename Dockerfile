FROM golang:1.18
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /usr/src/app
COPY . .
RUN go build -o ccblog .
EXPOSE 8888
ENTRYPOINT ["./ccblog --env prod"]