FROM golang:1.23.2

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update

# build go app
RUN go mod download
RUN go build -o crm-api ./cmd/main.go

RUN chmod +x crm-api

CMD ["./crm-api"]
