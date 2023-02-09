# syntax=docker/dockerfile:1

FROM golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main/* ./main/
COPY pb/* ./pb/
COPY pkg/* ./pkg/
RUN go build -o / ./main/main.go 

EXPOSE 8080

CMD [ "/main", "-grpc-addr", ":8080"]