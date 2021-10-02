# Mini-Madeuce

![](https://dockeri.co/image/nexusuw/mini-madeuce)

# WIP

```
touch db/mini-maduece
```

## todo

- caddyfile
  - basic logging
  - tor service
- blog post
- proper doc
- auto build
- ratelimiting?
- password parameter? ie /<stub>?p=password (url would be encrypted on sever + visiting url without pass would count against hits)
  - encrypt regardless if no password provided -> - basic encryption (https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/09.6.html) . encrypt url at rest, would provide protection if db were to be leaked

referrer+CSP is set in html to make sure that commited to code and reduces config requirement to deploy
