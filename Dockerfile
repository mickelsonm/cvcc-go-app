FROM golang:1.11 as builder

ENV GOPATH=/go/src
WORKDIR /go/src/cvcc
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod=vendor -o /cvcc

FROM alpine

WORKDIR /app/

COPY --from=builder /cvcc .
RUN ls
CMD [ "./cvcc" ]
