FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build -o /piscon_portal main.go


FROM alpine:3.14.0
WORKDIR /app

EXPOSE 4000

RUN apk add --no-cache --update ca-certificates imagemagick && \
  update-ca-certificates

COPY --from=build /piscon_portal /go/src/github.com/traPtitech/piscon-portal/.env  ./

ENTRYPOINT ./piscon_portal