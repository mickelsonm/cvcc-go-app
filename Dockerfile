FROM golang:1.11 as builder

ENV GOPATH=/go/src
WORKDIR /go/src/cvcc
COPY . .

RUN go build -mod=vendor -o /cvcc

FROM alpine

WORKDIR /app

COPY --from=builder /cvcc /app/cvcc
ENTRYPOINT ./cvcc
