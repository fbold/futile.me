### Migrations
Using `golang-migrate` for migrations: https://github.com/golang-migrate/migrate

Create migration:
```migrate create -ext sql -seq -dir internal/db/migrations```

Run migrations:
```
migrate -path internal/db/migrations -database "postgres://postgres:password@localhost:4001/futile?sslmode=disable" up
```
Note the `sslmode` param, that is to stop an "SSL not enabled" error.

