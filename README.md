# Authentication service

## Migrate guide
[Migrate repository](https://github.com/golang-migrate/migrate)
1. [Install migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
2. Export database URL
```shell
$ export POSTGRES_URL='postgres://user:password@localhost:5432/auth_db?sslmode=disable' 
```
2. Up migrations
```shell
$ migrate -database ${POSTGRES_URL} -path migrations up
```
3. Down migrations
```shell
$ migrate -database ${POSTGRES_URL} -path migrations down
```