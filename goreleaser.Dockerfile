FROM alpine:latest
LABEL org.opencontainers.image.authors=moulickaggarwal
RUN apk add git --no-cache
WORKDIR /
COPY k8s-cr-validator /
