FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build -o /piscon_portal main.go


FROM debian:latest
WORKDIR /app

EXPOSE 4000

COPY --from=build /piscon_portal /go/src/github.com/traPtitech/piscon-portal/.env  ./

ENTRYPOINT ./piscon_portal