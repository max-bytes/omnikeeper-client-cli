FROM alpine:3.14

RUN apk add \
    wget \
    curl \
    jq \
    libc6-compat # NOTE: required on alpine to run non-musl based binary, see https://stackoverflow.com/a/66974607
RUN wget https://github.com/max-bytes/omnikeeper-client-cli/releases/latest/download/okql
RUN chmod +x okql