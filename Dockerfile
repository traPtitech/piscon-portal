FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY ./go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /piscon_portal main.go


FROM alpine:3
WORKDIR /app

EXPOSE 4000

RUN apk update \
&& apk add ca-certificates \
&& update-ca-certificates

COPY --from=build /piscon_portal /go/src/github.com/traPtitech/piscon-portal/.env  ./

ENTRYPOINT ./piscon_portal
