# Go Rest Template

![Github Actions](https://github.com/kelvinator07/go-rest-template/actions/workflows/ci.yml/badge.svg)


# Gin
# JWT
# Redis for caching 
# Postgres / Mongodb

# API Structure
Main Server => Route => Handler => UseCase/Controller => Repository => Database

Route => routes/endpoints
Handler => payload validation, redis, caching
UseCase/Controller => status code, call to repositories

Follow instructions in `setup.txt` to add new entity
