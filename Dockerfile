FROM golang:latest

WORKDIR /app
COPY . .

ENV GO111MODULE=on

RUN go mod tidy
RUN go build -o app .

EXPOSE 8000

RUN chmod +x wait-for-it.sh
# Run the application after Mysql launch success
ENTRYPOINT ["./wait-for-it.sh", "mysql:3306", "--", "./app"]