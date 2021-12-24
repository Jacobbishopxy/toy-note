# Toy Note

Toy Note is a simple note taking app. It can be used to store notes and affiliated files, and tags can be added to notes as well. Furthermore, it provides searching functionality, such as searching by title, date or tags.

The purpose of this project is to learn:

- how to use GORM to deal with associations, such as many-to-many;
- how to use Mongo/GridFS to upload and download files;
- how to setup a logger with slicing functionality for both development and production;
- how to setup an openAPI documentation beneath Gin framework;
- ...

## Structure

The project structure is simple enough to understand:

```txt
toy-note
    ├── api
    │   ├── controller
    │   │   ├── note.go
    │   │   ├── query.go
    │   │   └── response.go
    |   |
    │   ├── entity
    │   │   ├── affiliate.entity.go
    │   │   ├── post.entity.go
    │   │   ├── tag.entity.go
    │   │   └── common.go
    |   |
    │   ├── persistence
    │   │   ├── mongo_test.go
    │   │   ├── mongo.go
    │   │   ├── postgres_test.go
    │   │   └── postgres.go
    |   |
    │   ├── service
    │   │   ├── note.service_test.go
    │   │   ├── note.service.go
    │   │   └── repository.go
    |   |
    │   ├── util
    │   │   ├── config_test.go
    │   │   └── config.go
    |   |
    │   └── api.go
    |
    ├── cmd
    │   └── app
    │       └── main.go
    |
    ├── docs
    │   ├── docs.go
    │   ├── swagger.json
    │   └── swagger.yaml
    |
    ├── env
    │   ├── dev.env
    │   └── prod.env
    |
    ├── logger
    │   └── logger.go
    |
    ├── go.mod
    ├── go.sum
    └── Makefile
```

## Routes

```txt
- [GET]         /get-tags
- [POST]        /save-tag
- [DELETE]      /delete-tag/:id
- [GET]         /get-posts
- [POST]        /save-post
- [DELETE]      /delete-post/:id
- [GET]         /download-file/:id
- [GET]         /search-posts-by-tags
- [GET]         /search-posts-by-title
- [GET]         /search-posts-by-time
```

Note:

- `save-post` only accepts `multipart/form-data`, this is due to the demand of uploading multiple files. Hence, the only way to pass `entity.Post` info is to convert it into a string, and put it into an extra text field (here we use `data`).

## Configuration

Please modify your configs under the `toy-note/env` folder. The only configs we have now is all about PostgreSQL and MongoDB.

## Development

```bash
make dev
```

## Production

```bash
make build
make prod
```

## OpenAPI

With a running server, please go visit [swagger docs](http://localhost:8080/docs/index.html) to view the API documentation.

Note:

- `make install-swag` can install swagger CLI, details see [Swaggo](github.com/swaggo/swag)
- `make swag-init` to generate the latest swagger docs.

## Dependencies

- [x] Gin: Web framework

- [x] Gorm: ORM

- [x] Mongo-driver: MongoDB driver

- [ ] Redis: Caching

- [x] Zap + Lumberjack: Logging

- [x] Swag: API Documentation

- [x] Viper: Project configuration
