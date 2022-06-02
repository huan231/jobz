FROM node:16.15-alpine AS builder

WORKDIR /usr/src/app

COPY package*.json ./
COPY . ./

RUN npm ci
RUN npm run build

FROM nginx:1.22-alpine

ARG COMMIT_HASH
ENV COMMIT_HASH=$COMMIT_HASH

COPY --from=builder /usr/src/app/dist /usr/share/nginx/html

COPY /nginx/templates /etc/nginx/templates

COPY envsubst-on-template.sh /docker-entrypoint.d
