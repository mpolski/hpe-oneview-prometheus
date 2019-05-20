FROM golang:latest as builder
ENV GOPATH /go/
WORKDIR /build/ 

RUN go get -v -d -tags "static-netgo" github.com/prometheus/client_golang/prometheus && \
    go get -v -d -tags "static-netgo" github.com/mpolski/oneview-golang/ov

COPY hpe-oneview-prometheus.go /build/

RUN CGO_ENABLED=0 GOOS=linux go build -v -a -tags "static netgo" -ldflags '-w' hpe-oneview-prometheus.go

FROM scratch
COPY --from=builder /build/ /app/
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/app/hpe-oneview-prometheus"]