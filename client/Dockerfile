FROM node:lts-alpine as build

WORKDIR /app

COPY . /app

RUN yarn install && yarn build || true


FROM nginx

COPY --from=build /app/dist /usr/share/nginx/html