FROM node:24-alpine AS build_web

RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app
COPY --chown=node:node apps/web/package*.json ./
USER node
RUN npm install

COPY --chown=node:node ./apps/web .

RUN npm run build

FROM golang:1.25-alpine AS build_server

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./apps/server ./apps/server

RUN go build -o /exorcist ./apps/server/cmd/exorcist


FROM golang:1.25-alpine AS exorcist

RUN apk update && apk add ffmpeg

COPY --from=build_web /home/node/app/dist /web

COPY --from=build_server /app/apps/server/migrations /migrations
COPY --from=build_server /exorcist /exorcist

EXPOSE ${PORT}

CMD [ "/exorcist" ]
