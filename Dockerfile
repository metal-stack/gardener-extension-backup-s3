FROM golang:1.22 AS builder

WORKDIR /go/src/github.com/metal-stack/gardener-extension-backup-s3

COPY . .

RUN make install

FROM alpine:3.19

WORKDIR /

COPY --from=builder /go/bin/gardener-extension-backup-s3 /gardener-extension-backup-s3

CMD ["/gardener-extension-backup-s3"]
