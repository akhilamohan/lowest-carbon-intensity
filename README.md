## Build and run the server


To build and run the server

```
docker-compose -f docker-compose.yaml up --build
```
Server will be listening on port 3000


## Sample request

Request:
```
curl --request GET \
-H 'Accept: application/json' \
--url "http://localhost:3000/slots?duration=121&contiguous=false"
```

Response:
```
[{"valid_from":"2025-01-18T03:30:00Z","valid_to":"2025-01-18T04:00:00Z","carbon":{"intensity":138}},{"valid_from":"2025-01-18T05:00:00Z","valid_to":"2025-01-18T05:30:00Z","carbon":{"intensity":139}},{"valid_from":"2025-01-18T04:00:00Z","valid_to":"2025-01-18T04:30:00Z","carbon":{"intensity":139}},{"valid_from":"2025-01-18T04:30:00Z","valid_to":"2025-01-18T05:00:00Z","carbon":{"intensity":139}},{"valid_from":"2025-01-18T05:30:00Z","valid_to":"2025-01-18T05:31:00Z","carbon":{"intensity":141}}]
```
