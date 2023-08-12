
FROM node:20-alpine AS dev
RUN apk add --no-cache libc6-compat
RUN apk update

WORKDIR /ui

COPY package.json package-lock.json ./
RUN  npm install

COPY . .
EXPOSE 5100

ENV PORT 5100

CMD ["npm", "start"]