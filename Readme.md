## About

For sake of simplicity the project uses `sqlite3` as db and a very basic ui. 

### Running the app with Docker

```shell
docker build . -t gymshark
docker run -p 8081:8081 gymshark
```

### Running the server locally

To run the server locally you need to have `go` and `sqlite3` installed.

1. Start server locally:

```shell
go run main.go
```

2. Check Health of the API by going to http://0.0.0.0:8081/health. You should see a response "OK"

3. The UI is available at http://0.0.0.0:8081/

## Run Tests

```shell
go test ./...
```

## Curl requests

### Save pack sizes

For simplicity all previous pack sizes are removed. In a production scenario we would need to update them by a unique
identifier.

```shell
curl -X POST -H "Content-Type: application/json" -d '{"packSizes":[250,500,1000,2000,5000]}' http://localhost:8081/api/save-pack-sizes
```
Response:
```json
{"success":true,"msg":"Created"}
```

### Get all pack sizes:

```shell
curl http://localhost:8081/api/pack-sizes
```

Response:
```json
{
  "packSizes": [
    111,
    222,
    333,
    444,
    555
  ]
}
```

### Calculate nr of packs

```shell
curl http://localhost:8081/api/packs/12001
```
Response:
```json
{
  "packs": [
    {
      "qty": 1,
      "items_per_unit": 2000
    },
    {
      "qty": 2,
      "items_per_unit": 5000
    },
    {
      "qty": 1,
      "items_per_unit": 250
    }
  ]
}
```