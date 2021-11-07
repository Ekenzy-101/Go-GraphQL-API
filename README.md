# GO GRAPHQL API

A GraphQL API using the [Clean Architecture by Uncle Bob`](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## TECH STACK

- Golang
- MongoDB
- PostgreSQL
- GraphQL ([graphql-go](https://github.com/graphql-go/graphql), any GraphQL library is fine as far you know the concepts of GraphQL)

## SETUP

- Copy the following environmental variables to a `.env` file from `example.env` and fill in your credentials

  ### WITH POSTGRESQL

  - Install a migrator e.g. [tern](https://github.com/jackc/tern) and create a `tern.conf` with the neccessary config
  - Run `make migrate`
  - Uncomment the postgres codes and comment the mongo codes in the following files
    - `repository/repository.go`
    - `main.go`
    - `service/common.go`

  ### WITH MONGODB

  - Nothing to do

- Run `make dev` to start the application

For any constructive feedback on this code, please open an issue Thanks!
