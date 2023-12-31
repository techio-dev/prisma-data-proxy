FROM golang:alpine3.18 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /data-proxy

FROM alpine:3.18
ARG DOCKER_METADATA_OUTPUT_VERSION
ENV VERSION $DOCKER_METADATA_OUTPUT_VERSION
COPY --from=builder /data-proxy /data-proxy
EXPOSE 8080
CMD [ "/data-proxy" ]
