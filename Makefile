 -database=postgres://arailym:artimona5@localho
st/anime?sslmode=disable

 curl -i -d "$BODY" localhost:4000/v1/animes
 curl -X PUT -d "$BODY" localhost:4000/v1/animes/3

 DELETE CURL:
 curl -X DELETE localhost:4000/v1/animes/2

 CREATE CURL:
 BODY='{"title":"Your name","year":2016,"runtime":"106 mins", "genres":["animation","fantasy","drama"]}'
  curl -d "$BODY" localhost:4000/v1/animes

 BODY='{"title":"Demon Slayer the Movie: Mugen Train","year":2020,"runtime":"117 mins", "genres":["animation","adventure","action"]}'
 curl -d "$BODY" localhost:4000/v1/animes

 BODY='{"title":"A Silent Voice","year":2016,"runtime":"130 mins", "genres":["animation","drama"]}'
 curl -d "$BODY" localhost:4000/v1/animes

 ADV UPDATE
 curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/animes/4

 go run ./cmd/api

 Rate limit - for i in {1.4.}; do curl http://localhost:4000/v1/healthcheck; done

 User CREATE
 BODY='{"name": "Alice Smith", "email": "alice@example.com", "password": "pa55word"}'
 curl -i -d "$BODY" localhost:4000/v1/users