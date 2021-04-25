## Book-API-Fiber
---
> API build using go-lang Fiber Framework

Refer: `https://docs.gofiber.io/api/middleware`

#### Usage:
```sh
docker-compose up -d
go run loader.go
```

#### API EndPoints:
```sh
NOTE: make sure to set Auth Header
HEADERS: {
    "Authorization": "Basic WVhCcFlubDJhWFpsYXdvPTphYzhkZWIyYjUzZGIzZGY4NDE3ZGU2ZmY5ZjI5NDVlYTYwMWQzYWIwOWM1MWY0OTEzMjM1NjNmOGE3M2U4M2Q2"
}

http://localhost:9091/book-api/v1
http://localhost:9091/book-api/v1/firstbook
http://localhost:9091/book-api/v1/latest/updated
http://localhost:9091/book-api/v1/latest/created
http://localhost:9091/book-api/v1/byauthor?name=Douglas Adams
http://localhost:9091/book-api/v1/byid/12
http://localhost:9091/book-api/v1/bookcount
http://localhost:9091/book-api/v1/bookid/range?from=5&to=13
```

#### Todo:
###### GET
- firstbook DONE
- lastbook DONE
- latestupdated DONE
- latestcreated DONE
- byid/:id DONE
- book/:?query in-progress
- bookcount DONE
- book/:?from=5&to=10 (5 to 10) DONE
- byauthor/:?name DONE
###### POST
- createbook CHECK
- updatebook CHECK
###### DELETE
- deletebook CHECK