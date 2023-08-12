
FROM node:20-alpine AS dev
RUN apk add --no-cache libc6-compat
RUN apk update

WORKDIR /ui

COPY package.json package-lock.json ./
RUN  npm install
ENV NODE_ENV=development
ENV CHOKIDAR_USEPOLLING=true
ENV WATCHPACK_POLLING=true
COPY . .
EXPOSE 5100

ENV PORT 5100

CMD ["npm", "start"]