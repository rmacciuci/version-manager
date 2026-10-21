FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /version-manager .

FROM alpine:3.22

LABEL "com.github.actions.name"="Version Manager"
LABEL "com.github.actions.description"="Manage versions of your project"
LABEL "com.github.actions.icon"="tag"
LABEL "com.github.actions.color"="blue"
LABEL "maintainer"="Ramiro Macciuci <ramimacciuci@gmail.com>"

RUN apk add --no-cache git
COPY --from=build /version-manager /usr/local/bin/version-manager

ENTRYPOINT ["/usr/local/bin/version-manager"]
