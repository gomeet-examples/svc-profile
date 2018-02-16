# golang builder
FROM golang:1.8.3-alpine3.6 AS builder
WORKDIR /go/src/github.com/gomeet-examples/svc-profile/
COPY . .
RUN apk add --no-cache --update git make protobuf protobuf-dev ca-certificates curl && \
     rm -rf /var/cache/apk/*
RUN rm -f /go/src/github.com/gomeet-examples/svc-profile/_build/svc-profile
RUN make tools-clean tools-sync tools
RUN make

# cf. https://hub.docker.com/r/gomeet/gomeet-builder/
# FROM gomeet/gomeet-builder:0.0.3 AS builder
# WORKDIR /go/src/github.com/gomeet-examples/svc-profile/
# COPY . .
# RUN rm -f /go/src/github.com/gomeet-examples/svc-profile/_build/svc-profile
# RUN make

# minimal image from scratch
FROM scratch
LABEL maintainer="Hugues Dubois <hugdubois@gmail.com>"
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/github.com/gomeet-examples/svc-profile/_build/svc-profile /svc-profile
EXPOSE 13000
ENTRYPOINT ["/svc-profile"]
CMD ["serve"]
