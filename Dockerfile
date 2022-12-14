FROM golang:latest
ENV TZ=Asia/Shanghai

WORKDIR /app
ADD go.mod .
ADD go.sum .

RUN export GO111MODULE=on && \
export GOPROXY=https://goproxy.cn && \
go mod download

COPY . .

RUN GO111MODULE=on go build -o app .

CMD ["./app", "-envTag", "docker"]