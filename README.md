# Gator
## A Blog Aggregator

### Prerequisites
 - Install GO (version: >=1.22)
 - Install Postgresql (version: >=16)

### Install
`go install github.com/ayushjaiswal22/gator@latest`

### Configuration
#### Create ~/.gatorconfig.json as below using db's connection string:
`{"db_url":"<db_connection_string>}`
#### For example:
`{"db_url":"postgres://<username>@localhost:5432/<dbname>?sslmode=disable"}`

#### Migrate the DB
```
cd $GOPATH/pkg/mod/github.com/user/repo@version/sql/schema/
goose postgres <db_url> up
```

##### In case migration fails:
```
goose postgres <db_url> down
goose postgres <db_url> up
```

