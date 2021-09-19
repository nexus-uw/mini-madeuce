FROM --platform=$BUILDPLATFORM golang:alpine as builder
WORKDIR /go/src/github.com/nexus-uw/mini-madeuc
COPY  . /go/src/github.com/nexus-uw/mini-madeuce/
RUN apk add --no-cache gcc musl-dev && go get github.com/mattn/go-sqlite3
RUN go get ./
RUN  go build -o /go/bin/mini-madeuce

FROM scratch
COPY --from=builder /go/bin/mini-madeuce /app/mini-madeuce
COPY ./template /app/template
COPY ./static /app/static

#ENV PORT 9090
# host where user accesses clearnet site
ENV HOST "http://localhost:9090"
# tor address for darknet site
ENV ONION "http://TODO.onion"

#sqlite db file
VOLUME /app/mini-madeuce.db

EXPOSE 9090
CMD ["/app/mini-madeuce"]
