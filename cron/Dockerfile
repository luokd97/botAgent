FROM golang:latest

WORKDIR /cron_task
ADD go.mod .
ADD go.sum .

RUN export GO111MODULE=on && \
export GOPROXY=https://goproxy.cn && \
go mod download

COPY . .

RUN GO111MODULE=on go build -o task .

CMD ["./task"]