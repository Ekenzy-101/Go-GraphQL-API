# GO GRAPHQL API

A GraphQL API using the [Clean Architecture by Uncle Bob`](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## TECH STACK

- Golang
- MongoDB
- PostgreSQL
- GraphQL ([graphql-go](https://github.com/graphql-go/graphql), any GraphQL library is fine as far you know the concepts of GraphQL)

## SETUP

- Copy the following environmental variables to a `.env` file from `example.env` and fill in your credentialsss
- Run `make migrate` if `DATABASE_TYPE=postgres`
- Run `make start-db` to start the databases 
- Run `make dev` to start the application

For any constructive feedback on this code, please open an issue Thanks!
