.DEFAULT_GOAL:= help
.PHONY:= run db-run-help

run: ## Run the application
	go run main.go

db-run: ## Run the database
	docker run --name mongodb -d -p 27017:27017 -v melisearch-db:/data/db  mongodb/mongodb-community-server

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep ^help -v | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
