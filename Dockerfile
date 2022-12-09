FROM golang:latest

WORKDIR /app
ADD go.mod .
ADD go.sum .

ENV GO111MODULE=on

RUN go mod tidy

COPY . .

RUN go build -o app .

EXPOSE 8000

RUN chmod +x wait-for-it.sh
# Run the application after Mysql launch success
#ENTRYPOINT ["./wait-for-it.sh", "mysql:3306", "--", "./app"]
ENTRYPOINT ["./app", "docker"]