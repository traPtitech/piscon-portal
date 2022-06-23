# ビルド環境
FROM node:14-alpine as build
WORKDIR /app
COPY . .
# https://github.com/docker/getting-started/issues/124
RUN apk --no-cache --virtual build-dependencies add \
        python3 \
        make \
        g++ \
&& npm install \
&& apk del build-dependencies
RUN npm run build

# 本番環境
FROM caddy:2.5.1-alpine
EXPOSE 80
COPY --from=build /app/dist /usr/share/caddy
