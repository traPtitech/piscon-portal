FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o piscon_portal main.go
RUN apk update \
&& apk add ca-certificates \
&& update-ca-certificates

FROM scratch
WORKDIR /app
EXPOSE 4000
COPY --from=build /go/src/github.com/traPtitech/piscon-portal/piscon_portal \
									/go/src/github.com/traPtitech/piscon-portal/.env  ./
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/piscon_portal"]
