FROM golang:1.22.1 as builder
COPY go.mod go.sum /go/src/github.com/oybek/choguuket/
WORKDIR /go/src/github.com/oybek/choguuket
RUN go mod download
COPY . /go/src/github.com/oybek/choguuket
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/choguuket github.com/oybek/choguuket

FROM alpine/curl
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/oybek/choguuket/build/choguuket /usr/bin/choguuket
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/choguuket"]
