include .env
export $(shell sed 's/=.*//' .env)

dev:
	@go run main.go

migrate:
	@tern migrate -m ./migrations
	
restart-db:
	# @sudo service postgresql restart
	@sudo service mongod restart

