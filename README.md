# Coopera

## Local Development Setup

You need **Docker** and **Docker Compose**.

### 1\. Initial Setup and Database Start

Run once to start PostgreSQL and apply all schema migrations:

```bash
docker-compose up -d postgres migrate
```

### 2\. Run the Go Application

After the database is running, use this variable on your host machine to connect the application:

```bash
DATABASE_URL="postgres://user:password@localhost:5432/database?sslmode=disable"
```

## Schema Management

### 1\. Apply New Migrations

To apply any pending schema changes:

```bash
docker-compose run --rm migrate up
```

### 2\. Rollback the Last Migration

To undo the most recent schema change:

```bash
docker-compose run --rm migrate down 1
```

### 3\. Stop All Services

To stop and remove all running containers and networks:

```bash
docker-compose down
```
