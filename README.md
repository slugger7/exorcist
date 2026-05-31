# Exorcist

**!!! UNDER CONSTRUCTION !!!**

A project to manage your media and reduce duplicates.

Link people to videos and photos

Add tags to media

## Getting started

### Prerequisites

- [Docker](https://www.docker.com/)

### Running application

- Copy [.env.example](.env.example) to [.env](.env)
  - Alter varibales in [.env](.env) to suit your needs
- In [docker-compose.yml](docker-compose.yml) mount the volumes with your media to the server
- Bring the environment up

```bash
docker compose up
```

- In your browser go to [http://localhost:3366](http://localhost:3366) (unless configured otherwise)
- Log in with the username `admin` and password `admin`
  - In the top right when you head to the admin profile you can change the password
- Add a library
- Add a library path to that library
- Scan the library path

If you do not require pgAdmin you can modify your [docker-compose.yml](./docker-compose.yml) and remove 
the pgadmin container and its volume

## Development

### Prerequisites

- [Go](https://go.dev/)
- [Node](https://nodejs.org/en)
- [ffmpeg](https://ffmpeg.org/)
- [Docker](https://www.docker.com/)
- [psql](https://www.postgresql.org/docs/current/app-psql.html)
- [tygo](https://github.com/gzuidhof/tygo)
  - Mac:
    - `brew install libpq`
    - `brew link --force libpq`
  - Ubuntu:
    - `sudo apt update`
    - `sudo apt install postgresql-client`

### Running

## Database

- `docker compose up db -d`
- There is a VS Code plugin suggestion for postgres that you could use to manage the database
- You could also add a pg_admin container into the docker compose

## pgAdmin

In order to control and manage your database you could use the pgadmin container.

- `docker compose up pgadmin -d`
- <http://localhost:5444>
- Username and password are in the [.env](./.env) file `PGADMIN_DEFAULT_EMAIL` and `PGADMIN_DEFAULT_PASSWORD` 
- Right click on the serves
- Select register
- Name it
- *Host:* `db`
  uses docker dns to resolve
- *Username*: `exorcist` (found in [.env](./.env) `DATABASE_USER`)
- *Password*: `some-super-secret` (found in [.env](./.env) `DATABASE_PASSWORD`)
- Optional: save the password
- Click save
- You shoud be able to find the new database connection in the servers list

## Server

In order to do some development you can run the server by itself.

From the root directory:

- `go run ./apps/server/cmd/exorcist`

There are also some configuration files for some editors but do not expect them to work or be up to date.

## Web

In order to run the web you will need to do the following

- `cd ./apps/web`
- `cp .env.example .env`
- modify [.env](./.env) with the parameters you desire
- `npm run dev`
