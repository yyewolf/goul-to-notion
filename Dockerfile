FROM golang:1.17.3-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN export CGO_ENABLED=0 && go build -o /notionUpdater

##
## Deploy
##
FROM alpine:3.14

WORKDIR /

RUN mkdir -p /data
RUN chmod 777 /data

COPY --from=build /notionUpdater /notionUpdater

ENTRYPOINT ["/notionUpdater"]