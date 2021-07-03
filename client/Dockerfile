FROM node:14.17 as build
WORKDIR /app
COPY . /app
RUN npm install && npm run build

FROM nginx
EXPOSE 80
COPY --from=build /app/dist /usr/share/nginx/html