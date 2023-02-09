FROM golang:latest AS builder

RUN apt-get update && apt-get install -y ca-certificates openssl git tzdata

ARG cert_location=/usr/local/share/ca-certificates

COPY /certs/xkcd.crt /usr/local/share/ca-certificates/xkcd.crt
# Update certificates
RUN update-ca-certificates

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIT_SSL_NO_VERIFY=1

WORKDIR /go/src

COPY go.mod .

RUN go mod download

COPY . .

RUN go build main.go

RUN go build notconcurrent.go

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# This line will copy all certificates to final image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src .

# Exec the binary file produced by `go build` according to parameters passed to `docker run`
ENTRYPOINT ["./main"]
