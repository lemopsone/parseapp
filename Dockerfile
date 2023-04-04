FROM golang:latest as builder
COPY go.mod go.sum /go/src/github.com/lemopsone/parseapp/
WORKDIR /go/src/github.com/lemopsone/parseapp
RUN go mod download
COPY . /go/src/github.com/lemopsone/parseapp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./app github.com/lemopsone/parseapp

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/lemopsone/parseapp /usr/bin/parseapp
EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/parseapp/app"]