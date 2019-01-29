FROM alpine AS base
RUN apk add --no-cache curl wget

FROM golang:1.11 AS go-builder
WORKDIR /go/app
COPY . /go/app
RUN go get github.com/satori/go.uuid
RUN go get github.com/lib/pq
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/app/book-master-api /go/app/src/book-master-api.go

FROM base
COPY --from=go-builder /go/app/book-master-api /book-master-api
CMD ["/book-master-api"]