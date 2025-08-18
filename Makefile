include .env
export $(shell sed 's/=.*//' .env)

deploy-service:
	@kubectl create secret generic $(NAME) --from-env-file $(NAME).env
	@kubectl apply -f k8s/$(NAME).yaml

deploy-ingress:	
	@gcloud compute addresses create backend --global
	@kubectl apply -f k8s/ingress.yaml

delete-ingress:
	@gcloud compute addresses delete backend --global
	@kubectl delete -f k8s/ingress.yaml

delete-service:
	@kubectl delete secret $(NAME)
	@kubectl delete -f k8s/$(NAME).yaml

dev:
	@go run main.go

migrate:
	@go install github.com/jackc/tern/v2@latest
	@tern migrate -m ./migrations --conn-string $(DATABASE_URL)
	
prod:
	@docker compose up -d api

start-db:
	@docker compose up -d cache $(DATABASE_TYPE) 

stop-db:
	@docker compose stop cache $(DATABASE_TYPE)


