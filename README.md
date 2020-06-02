# Desafio Backend - Star Wars

## Running the project
### Requirements
- [Golang](https://golang.org)
- [MongoDB](https://www.mongodb.com)
- Internet connection (to extract informations of [swapi.dev](https://swapi.dev))

### Steps
1. Clone the repository: `git clone git@github.com:mauriciolopes/desafio-b2w`
2. Go to `cmd/planet-api` folder and run: `go install`
3. Run: `$GOPATH/bin/planet-api`

## Running the project with docker
### Requirements
- [Docker](https://www.docker.com)
- [Docker Compose](https://www.docker.com)
- Internet connection (to extract informations of [swapi.dev](https://swapi.dev))

### Steps
1. Clone the repository: `git clone git@github.com:mauriciolopes/desafio-b2w`
2. Inside the root project folder: `docker-compose up --build`

## Configuration
### Environment variables
- `DB_URI`: configure database server (default `mongodb://localhost:27017`)
- `DB_NAME`: configure database name (default `starwars`)
- `HTTP_ADDR`: configure http server host and port (default `:8080`)

### Binary flags
These configurations have priority over environment variables.
- `-db-uri`: configure database server (default `mongodb://localhost:27017`)
- `-db-name`: configure database name (default `starwars`)
- `-http-addr`: configure http server host and port (default `:8080`)

### Examples

```bash
DB_URI=mongodb://localhost:27017 DB_NAME=starwars planet-api
```

```bash
HTTP_ADDR=localhost:8088 planet-api -db-uri=mongodb://localhost:27017
```

## How to use

The content type of requests and responses is `application/json`.

### Create a planet

`POST /planets`

**Request**

| Field | Type | Description | Rules |
|---|---|---|---|
| name | string | Name of the new planet | required, unique |
| climate | string | Climate of the new planet | required |
| terrain | string | Terrain of the new planet | required |

**Response example**

HTTP Status: 201 Created

```json
{
    "id": "5ed656285d9f32274f847473",
    "name": "Hoth",
    "climate": "frozen",
    "terrain": "tundra, ice caves, mountain ranges",
    "filmAppearances": 1
}
```

### List all planets

`GET /planets`

**Response example**

HTTP Status: 200 OK

```json
[
  {
    "id": "5ed656285d9f32274f847473",
    "name": "Hoth",
    "climate": "frozen",
    "terrain": "tundra, ice caves, mountain ranges",
    "filmAppearances": 1
  },
  {
    "id": "5ed659005d9f32274f847474",
    "name": "Naboo",
    "climate": "temperate",
    "terrain": "grassy hills, swamps, forests, moutains",
    "filmAppearances": 4
  },
  {
    "id": "5ed659415d9f32274f847475",
    "name": "Kamino",
    "climate": "temperate",
    "terrain": "ocean",
    "filmAppearances": 1
  }
]
```

### Get a planet

`GET /planets/:id`

**Response example**

```json
{
  "id": "5ed659005d9f32274f847474",
  "name": "Naboo",
  "climate": "temperate",
  "terrain": "grassy hills, swamps, forests, moutains",
  "filmAppearances": 4
}
```

### Get a planet by name

`GET /planets?name=:name`

**Response example**

```json
{
  "id": "5ed659005d9f32274f847474",
  "name": "Naboo",
  "climate": "temperate",
  "terrain": "grassy hills, swamps, forests, moutains",
  "filmAppearances": 4
}
```

### Delete a planet

`DELETE /planets/:id`

**Response example**

HTTP Status: 200 OK
