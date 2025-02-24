# Exorcist

Similar to [ghost](https://github.com/slugger7/ghost-media) and [poltergeist](https://github.com/slugger7/poltergeist) this one is written in golang

## Getting started

- Install Go
- Install Docker
- Install psql
  - Mac:
    - `brew install libpq`
    - `brew link --force libpq`
- Install ffmpeg
- Copy `templte.env` -> `.env` and fill in the details
- `docker compose up -d` to start the database
- `make run` to start the application

## Tools

### Api

[Gin](https://go.dev/doc/tutorial/web-service-gin)

#### Authentication

[example](https://github.com/depado/gin-auth-example/tree/main)

### Database

#### Entity Relation Diagram

![entity_relation_diagram](./diagrams/out/entity_relation_diagram.d2.svg)

This application is using a database first approach with a sql builder instead of an ORM.
It utilizes [Jet](https://github.com/go-jet/jet)

#### Migrations

We utilize [golang-migrate](https://github.com/golang-migrate/migrate) for our migrations

Migrations run automatically when the application starts up.
It is recommended to install the cli tool for running migrations. This will allow you to run migrations from the command line without having to start up the application. [CLI migration](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
It is also recommended that you install the [Jet CLI tool](https://github.com/go-jet/jet?tab=readme-ov-file#prerequisites) in order to update the models of the application

- Create migration: `./scripts/create-migration.sh <migration-name>`
- Run migrations: `./scripts/run-migrations.sh`
- Undo a migration: `./scripts/undo-migration.sh`
- Update models: `./scripts/update-models.sh`

The usual workflow would to add a migration would be:

1. Create a migration
1. Run the migrations
1. Update the models

### FFMpeg stuff

[Wrapper](https://github.com/u2takey/ffmpeg-go)

### Environment

[dotenv](https://github.com/joho/godotenv)

### Troubleshooting

- If using zsh remember to add the following to your .zshrc

  ```bash
  export GOPATH=$HOME/go  
  export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
  ```

### Diagrams

Generating of diagrams is done by utilizisg [d2](https://d2lang.com/)
In order to generate diagrams you will need to install the cli tool

### Mocks

Currently there is a [mocks](./internal/mocks) that have mocks in that have been made by us. This is *not* to be used any more as it is too much effort to manually update these mocks.
As these are not being used any more the methods that are not being used any more should panic to prevent anyone else from using these mocks.

We are using a mock generating library to prevent us having to create these mocks by hand. To get this working:

- Install the cli tool [mockgen](https://github.com/uber-go/mock)
- When an interface has changed you should be able to run `make mocks` and the mocks will be regenerated for you
