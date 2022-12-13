#
# tiny go builder
#

FROM ubuntu:20.04

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
	golang-go build-essential ca-certificates

WORKDIR /host

CMD go get && make build
