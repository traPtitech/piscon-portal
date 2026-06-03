# ビルド環境
FROM node:14-alpine@sha256:434215b487a329c9e867202ff89e704d3a75e554822e07f3e0c0f9e606121b33 as build
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
FROM caddy:2.5.1-alpine@sha256:0033b34d2df3fe0bf94088c36e7d722ceca1b38cbdd49c08b2c10b9f9aa58912
EXPOSE 80
COPY --from=build /app/dist /usr/share/caddy
