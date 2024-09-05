
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

# Formas de programar MongoDB en GO

## Tener en cuenta cuando hacemos queries 

###  Usar projection cuando sea necesario:
- Projection es una palabra cheta que decidieron ponerle a algo muy simple (como todo en informatica)
- Simplemente es un filtro de partes del documento queremos. 
	- `{"field", 1}` si lo queremos
	- `{"field", 0}` si no lo queremos. No es necesario especificar los que no queremos, pero es buena practica

Entonces si no somos fiacas deberiamos usar projection siempre. Es hasta mas performante en la mayoria de los casos. Peor bueno si tenemos un documento que es pequeño vamos para adelante sin projection

La cosa es que si no usamos es mucho mas comodo porque en go el decode lo va a hacer igual:

Caso sin projection:
```go
	var result struct {
		A string `bson:"A"`
	}

	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return "", err
	}

	return result.Word, nil
```

Caso con projection:
```go
	onlyFieldA := bson.M{"A", 1}

	var result struct {
		A string `json:"A"`
	}

	if err := col.FindOne(ctx, filter, options.FindOne().SetProjection(onlyFieldA)).Decode(&result); err != nil {
		return "", err
	}

	return result.Word, nil
```

En el primer caso solo le esta llegando el field A. 

En el segundo caso le esta llegando todo pero como que nos enteramos porque estamos haciendo el decode a la misma estructura. De todos formas en este caso le llega toda la estructura, esto puedo incluir un monton de trafico de red de mas y de procesamiento innecesario. 

# Notas
- Para asegurar que no se repita "word" en MongoDB creamos un index `unique_word_index`

