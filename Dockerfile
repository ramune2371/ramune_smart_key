FROM golang:1.21.3-alpine3.18

COPY ./linebot ./work

WORKDIR /go/work

RUN go mod download
RUN go build -o main /go/work/main.go

CMD /go/work/main
