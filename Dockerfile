FROM golang:1.14-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build -o /piscon_portal main.go


FROM alpine:3.12.0
WORKDIR /app

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

EXPOSE 4000

COPY --from=build /piscon_portal /go/src/github.com/traPtitech/piscon-portal/.env  ./

ENTRYPOINT ./piscon_portal