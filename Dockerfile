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

COPY ./server/go.mod .
COPY ./server/go.sum .
RUN go mod download

COPY ./server .

RUN go build -o /exorcist ./cmd/exorcist


FROM golang:1.24-alpine AS exorcist

RUN apk update && apk add ffmpeg

COPY --from=build_web /home/node/app/dist /web

COPY --from=build_server /exorcist /exorcist

EXPOSE ${PORT}

CMD [ "/exorcist" ]