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
	@tern migrate -m ./migrations
	
prod:
	@docker compose up -d api

start-db:
	@docker compose up -d cache db 

stop-db:
	@docker compose stop cache db


