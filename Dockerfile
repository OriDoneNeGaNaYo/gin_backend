FROM  golang:1.19

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN go build -o app Main.go

EXPOSE 8080

CMD ["./app"]