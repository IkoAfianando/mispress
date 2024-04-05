FROM golang:1.21.9-alpine3.19

COPY go.mod go.sum /go/src/github.com/IkoAfianando/mispress/

WORKDIR /go/src/github.com/IkoAfianando/mispress
RUN go mod download
COPY . /go/src/github.com/IkoAfianando/mispress

RUN go build -o /usr/bin/mispress cmd/main.go

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/mispress"]