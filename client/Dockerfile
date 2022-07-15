# ビルド環境
FROM node:14-alpine as build
WORKDIR /app
# https://github.com/docker/getting-started/issues/124
RUN apk --no-cache add \
      python3 \
      make \
      g++
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm install
COPY . .
RUN npm run build

# 本番環境
FROM caddy:2.5.1-alpine
EXPOSE 80
COPY --from=build /app/dist /usr/share/caddy
