FROM golang:1.22 AS builder

WORKDIR /go/src/github.com/metal-stack/gardener-extension-backup-s3

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make install

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /go/bin/gardener-extension-backup-s3 /gardener-extension-backup-s3

CMD ["/gardener-extension-backup-s3"]
