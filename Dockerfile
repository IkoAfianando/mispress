FROM golang:1.21.9-alpine3.19

COPY go.mod go.sum /go/src/gitlab.com/repo/mispress/
WORKDIR /go/src/gitlab.com/repo/mispress
RUN go mod download
COPY . /go/src/gitlab.com/repo/mispress
RUN go build -o /usr/bin/mispress gitlab.com/repo/mispress

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/mispress"]