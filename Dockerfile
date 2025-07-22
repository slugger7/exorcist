FROM node:24-alpine AS build_web

RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app
COPY --chown=node:node web/package*.json ./
USER node
RUN npm install

COPY --chown=node:node ./web .

RUN npm run build

FROM golang:1.24-alpine AS build_server

WORKDIR /app

RUN apk update && apk add ffmpeg
COPY ./server/go.mod .
COPY ./server/go.sum .
RUN go mod download

COPY ./server .

COPY --from=build_web /home/node/app/dist /web

RUN go build -o /exorcist ./cmd/exorcist

EXPOSE ${PORT}

CMD [ "/exorcist" ]