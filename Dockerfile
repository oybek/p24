FROM golang:1.22.1 as builder
COPY go.mod go.sum /go/src/github.com/oybek/p24/
WORKDIR /go/src/github.com/oybek/p24
RUN go mod download
COPY . /go/src/github.com/oybek/p24
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/p24 github.com/oybek/p24

FROM alpine/curl
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/oybek/p24/build/p24 /usr/bin/p24
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/p24"]
