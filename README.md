# div-dash

## Development

-   `./web` React Frontend

## DB Migrations

```
$ goose -dir sql/migrations postgres "user=postgres password=postgres dbname=postgres sslmode=disable" status

```

## Naming Conventions

### Database

#### Constraints

-   {tablename}_{columnname(s)}_{suffix} // https://til.cybertec-postgresql.com/post/2019-09-02-Postgres-Constraint-Naming-Convention/
