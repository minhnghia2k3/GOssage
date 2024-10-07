Tech stack

- Go 1.23.1 standard library
- Docker
- Postgres running on Docker
- Using trigram search reduce query time

```
1. Query 200,000 rows without indexes
   Planning Time: 0.468 ms
   Execution Time: 60.410 ms
   (6 rows)
2. With indexes
   Planning Time: 0.516 ms
   Execution Time: 62.005 ms
3. With Postgres Gin indexes
   Planning Time: 0.602 ms
   Execution Time: 2.897 ms
```

- Swagger for docs
- Golang migrate for migrations
- Setting up automation workflows (audit,release,versioning...)




