# PostgreSQL Database Handler


This part needs bin-data:

```bash
go get -u github.com/go-bindata/go-bindata/...
```

To generate migrations install `migrate` using the tutorial at https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Then:

```bash
migrate create -ext sql -dir migrations -seq name_of_migration
```

It will be added automagically with go-generate to bindata
