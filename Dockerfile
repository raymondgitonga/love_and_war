FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . /app

RUN go build -o /love_and_war
EXPOSE 8081
CMD ["/love_and_war"]





