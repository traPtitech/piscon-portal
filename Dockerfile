FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/traPtitech/piscon-portal
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o piscon_portal main.go

FROM gcr.io/distroless/base
WORKDIR /app
EXPOSE 4000
COPY --from=build /go/src/github.com/traPtitech/piscon-portal/piscon_portal \
									/go/src/github.com/traPtitech/piscon-portal/.env  ./
ENTRYPOINT ["/app/piscon_portal"]
