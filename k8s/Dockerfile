FROM alpine:latest

WORKDIR /go

COPY --chmod=0755 build_artifact_bin hdb-bin
ADD templates/ templates/
COPY fixtures/testconfig.yml config.yml

ENTRYPOINT ["/go/hdb-bin"]
