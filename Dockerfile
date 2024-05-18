FROM alpine:3.12

LABEL "com.github.actions.name"="Version Manager"
LABEL "com.github.actions.description"="Manage versions of your project"
LABEL "com.github.actions.icon"="tag"
LABEL "com.github.actions.color"="blue"

LABEL "mantainer" = "Ramiro Macciuci <ramimacciuci@gmail.com>"

RUN apk add --no-cache git bash

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]
