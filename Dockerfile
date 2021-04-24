FROM golang:latest
ENV GOPROXY=direct
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o main ./services/loader/main.go
CMD ["/app/main"]
