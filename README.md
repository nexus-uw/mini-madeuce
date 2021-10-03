# Mini-Madeuce

![](https://dockeri.co/image/nexusuw/mini-madeuce)

# what

selfhosted url shortener written in go + run using docker

# why

mokin-token.ramsay.xyz notes are too long to manually type into a browser nav bar

# WIP

```
touch db/mini-maduece
go get
go build && PORT=9090 HOST=localhost:9090 ./mini-madeuce
```

goto [http://localhost:9090](http://localhost:9090)

## todo

- proper doc
- ratelimiting?
- password parameter? ie /<stub>?p=password (url would be encrypted on sever + visiting url without pass would count against hits)
  - encrypt regardless if no password provided -> - basic encryption (https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/09.6.html) . encrypt url at rest, would provide protection if db were to be leaked
  - too slow right now

referrer+CSP is set in html to make sure that committed to code and reduces config requirement to deploy

onion url is handled by caddy currently (dont want to make it configurable from code today)
