# Gossage

**Gossage** is a modern social platform designed to foster user interaction through follows, post creation, and comment
engagement. Built with Go (Golang), it leverages PostgreSQL for data storage, Redis for caching, and Docker for
streamlined deployment.

## Tech Stack

- Go (1.23.1): Utilizing the Go standard library for efficient, scalable backend development.
- PostgreSQL: Data persistence, running in a Dockerized environment, with performance optimization using trigram
  indexes.
- Redis: Caching for high-speed data retrieval and reduced database load.
- Docker: Containerized deployment ensuring consistency across environments.

## Get started

1. **Copy environment configuration**

        $ cp .env.example .env

2. **Build Docker containers**

        $ make compose.up

3. **Run database migrations**

       $ make migrate.up

4. Seed sample data

       $ make seed

5. Start the web application

       $ cd web
       $ npm run dev

6. Start the Go server

         $ go run ./cmd/api

## Performance Optimization

Trigram search is employed to significantly reduce query execution times for large datasets:

```
1. Query 200,000 rows without indexes:
    - Planning Time: 0.468 ms
    - Execution Time: 60.410 ms

2. With traditional indexes:
    - Planning Time: 0.516 ms
    - Execution Time: 62.005 ms

3. With PostgreSQL GIN indexes:
    - Planning Time: 0.602 ms
    - Execution Time: 2.897 ms
```

## Additional Features

- Swagger: API documentation and testing interface.
- Golang Migrate: Schema migrations for database versioning.
- CI/CD Automation: Automated workflows for auditing, versioning, and release management.