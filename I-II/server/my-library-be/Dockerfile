FROM golang:1.21.3-alpine as builder
WORKDIR /root/
RUN apk add --update --no-cache make

RUN mkdir -p /root/my-library
WORKDIR /root/my-library
COPY . .
RUN make build

FROM alpine:3.18
WORKDIR /root/

COPY --from=builder /root/my-library/bin/ .

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.8.0/wait wait
RUN chmod +x wait

ENV LIBRARY_MONGO=mongodb

EXPOSE 8080
CMD ./wait && ./my-library
