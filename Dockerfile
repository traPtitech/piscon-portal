FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o piscon_portal main.go

FROM ubuntu:22.04
WORKDIR /app
EXPOSE 4000
RUN apt update \
&& apt install -y ca-certificates \
&& update-ca-certificates
COPY --from=build /go/src/github.com/traPtitech/piscon-portal/piscon_portal \
									/go/src/github.com/traPtitech/piscon-portal/.env  ./
ENTRYPOINT ["/app/piscon_portal"]
