
# melisearch
## MongoDB
### Set Up
La forma mas fácil es usar docker para correrlo localmente con este comando:
```bash
docker run --name mongodb-melisearch -d -p 27017:27017 -v melisearch-db:/data/db  mongodb/mongodb-community-server
```

En el makefile tenemos dos reglas para esto
```Makefile
db-init: ## Init the database
	docker rm mongodb -f
	docker run --name mongodb-melisearch -d -p 27017:27017 -v melisearch-db:/data/db  mongodb/mongodb-community-server

db-run: ## Run the Database
	docker start mongodb-melisearch
```

Si ya tenemos corriendo la base podemos correr el programa con `make run`

### [From MySQL to MongoDB](https://www.mongodb.com/docs/manual/reference/sql-comparison/)

| MySQl    | MongoDB    |
| -------- | ---------- |
| Database | Database   |
| Table    | Collection |
| Row      | Document   |
| Column   | Field      |

# Notas
- Para asegurar que no se repita "word" en MongoDB creamos un index `unique_word_index`

