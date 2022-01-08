FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /go

COPY build_artifact_bin hdb-bin
ADD templates/ templates/
COPY fixtures/testconfig.yml config.yml

RUN chmod 755 /go/hdb-bin
ENTRYPOINT ["/go/hdb-bin"]
