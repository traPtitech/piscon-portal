FROM node:14-alpine as build

WORKDIR /app

COPY . /app

RUN yarn install && yarn build || true


FROM nginx:1.21.0

COPY --from=build /app/dist /usr/share/nginx/html