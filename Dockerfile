FROM golang:1.17.3-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go build -o /notionUpdater

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /notionUpdater /notionUpdater

CMD ["/notionUpdater"]