# Tasks for NIX Education

`cmd/echo` -- REST API server based on ECHO-framework.

`cmd/echo/api` -- api code with tests.

`cmd/echo-webserver` -- HTTP-server based on ECHO-framework with OAuth support.

GORM and SQLite are used for storage.

## echo-webserver

Use `go run ./cmd/echo-webserver/` or `run-webserver.sh` to run webserver.

`-conf` flag is used to provide path to JSON configuration file. Default path is `./conf.json`.

`conf.json` example:

```json
{
    "GoogleOAuth": {
        "ClientID": "id",
        "ClientSecret": "secret"
    },
    "FacebookOAuth": {
        "ClientID": "id",
        "ClientSecret": "secret"
    },
    "SessionsKey": "SESSIONS_KEY",
    "DSN": "storage.db?_foreign_keys=ON",
    "Port": "8080"
}
```

`_foreign_keys=ON` flag has to be provided for SQLite support of foreign keys which are used to relate users with posts and posts with comments.
