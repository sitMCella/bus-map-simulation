FROM node:lts-alpine as build
WORKDIR /app
COPY ./frontend/package*.json ./
RUN npm install
COPY ./frontend .
RUN npm run build

FROM nginx
WORKDIR /app
COPY --from=build /app/dist /app
COPY ./frontend/nginx.conf /etc/nginx/nginx.conf
EXPOSE 80