
# Maps to Waze API

This is the backend of the maps-to-waze-app.

The goal of this service is to provide a set of API endpoints that allow you to convert Google Maps links into Waze links.
## API Reference

#### Convert a Google Maps URL into a Waze URL

```http
  POST /convertURL
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `url` | `string` | **Required**. The google maps URL |

#### Get the static map

```http
  GET /staticMap?lat=&lon=
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `lat`      | `string` | **Required**. Latitude |
| `lon`      | `string` | **Required**. Longitude |

#### Get details about a place

```http
  GET /placeDetails?lat=&lon=
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `lat`      | `string` | **Required**. Latitude |
| `lon`      | `string` | **Required**. Longitude |


## Run Locally

Clone the project

```bash
  git clone https://github.com/daniel203/maps-to-waze-api/
```

Go to the project directory

```bash
  cd maps-to-waze-api
```

Edit the `env` file and rename it to `.env`

```bash
  mv env .env
```

Start the server

```bash
  go run .
```

