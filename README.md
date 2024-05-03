# Go Rest Template

![Github Actions](https://github.com/kelvinator07/go-rest-template/actions/workflows/workflow.yaml/badge.svg)


### Gin, JWT, Redis for caching, Postgres / Mongodb

## API Structure
Main Server => Route => Handler => UseCase/Controller => Repository => Database

Route => routes/endpoints

Handler => payload validation, redis, caching

UseCase/Controller => status code, call to repositories

## Run App
Setup the database `make db-setup`
Start the server `make serve`

Follow instructions in `setup.txt` to add new entity

