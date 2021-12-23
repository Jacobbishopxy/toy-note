# Toy Note

Toy Note is a simple note taking app. It can be used to store notes and affiliated files, and tags can be added to notes as well. Furthermore, it provides searching functionality, such as searching by title, date or tags.

The purpose of this project is to learn how to use GORM and Mongo/GridFS to build a simple project.

## Configuration

Please modify your configs under the `toy-note/env` folder. The only configs we have now is all about PostgreSQL and MongoDB.

## Development

```bash
make dev
```

This will start the development server. Please visit [swagger docs](http://localhost:8080/docs/index.html) to view the API documentation.

## Production

```bash
make build
make prod
```

## Dependencies

- [x] Gin: Web framework

- [x] Gorm: ORM

- [x] Mongo-driver: MongoDB driver

- [ ] Redis: Caching

- [x] Zap + Lumberjack: Logging

- [x] Swag: API Documentation

- [x] Viper: Project configuration
