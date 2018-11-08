FROM golang:1.11

ENV GOPATH=/go/src
WORKDIR /go/src/cvcc
COPY . .

RUN go build -mod=vendor -o cvcc

EXPOSE 8080

CMD ["./cvcc"]
