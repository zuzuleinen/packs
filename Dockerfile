FROM golang:1.22

RUN apt-get update && \
    apt-get install -y sqlite3 libsqlite3-dev

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main
RUN chmod +x main

EXPOSE 8081

CMD ["./main"]
