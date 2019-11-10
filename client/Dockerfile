FROM node:lts-alpine

WORKDIR /app

COPY . /app

EXPOSE 8080

RUN yarn global add http-server

RUN yarn install && yarn build || true

CMD http-server ./dist
