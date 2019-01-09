# Availabe routes (9.1.19)

| Method   | Route                | Query    | Description                                                                       |
|----------|----------------------|----------|-----------------------------------------------------------------------------------|
| GET      | /                    |          | Get status of webservice                                                          |
| GET      | /content             | lang, id | Get content as a JSON Response. Query Parameters are optional (will return all)   |
| POST     | /content             |          | Add/update one instance. Id and Language will be queried acording to posted json. |
| PUT      | /content             |          | same as POST                                                                      |
| DELETE   | /content/{lang}/{id} |          | Deletes instance from DB                                                          |
| GET/POST | /init                |          | Creates Dummy data to initialize DB (test purpose, will be removed)               |
| GET/POST | /reset               |          | Reset DB (test purpose, will be removed)                                          |

For route `/content?lang='es'` for example we get

```json
[
    {
        "question": "test 1",
        "answer": "test1 answer",
        "id": 4,
        "lang": "es",
        "category": "work",
        "created_at": "2019-01-09T14:30:40.858+01:00"
    },
    {
        "question": "test 1",
        "answer": "test1 answer",
        "id": 9,
        "lang": "es",
        "category": "work",
        "created_at": "2019-01-09T14:30:45.859+01:00"
    }
]
```

as an answer. 

