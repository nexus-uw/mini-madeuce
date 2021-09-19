#FROM golang:alpine as builder
FROM --platform=$BUILDPLATFORM golang:alpine as builder
RUN apk add gcc musl-dev
WORKDIR /go/src/github.com/nexus-uw/mini-madeuce
COPY  . /go/src/github.com/nexus-uw/mini-madeuce
RUN go get ./
RUN  go build -o /go/bin/mini-madeuce
RUN touch mini-madeuce.db

FROM alpine
#FROM scratch (something internall wants to use sh (maybe the sqlite lib)
WORKDIR /app
COPY --from=builder /go/bin/mini-madeuce /app/mini-madeuce
COPY --from=builder /go/src/github.com/nexus-uw/mini-madeuce/mini-madeuce.db  /app/db/mini-madeuce.db
COPY ./template /app/template
COPY ./static /app/static
#ENV PORT 9090
# host where user accesses clearnet site
ENV HOST "http://localhost:9090"
# tor address for darknet site
ENV ONION "http://TODO.onion"

#sqlite db file
VOLUME /app/db

EXPOSE 9090

#CMD "/app/mini-madeuce"
ENTRYPOINT /app/mini-madeuce
