# ビルド環境
FROM node:lts-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# 本番環境
FROM caddy:latest
EXPOSE 80
COPY dev/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /app/dist /usr/share/caddy