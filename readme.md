## Lotus backend

### Dependencies
 - Go v1.11
 - Postgres

### Build
With go v1.11 `go build` (thanks to modules)

Make sure you have a postgres DB and user setup:
 - Create a postgres user `createuser -P lotus`
 - Create a postgres DB `createdb -O lotus lotus`
 - Link admin to the frontend path `ln -s ~/Documents/node/lotus_frontend/ admin`

Create an env file like the following:
