# Toy Note

Toy Note is a simple note taking app. It can be used to store notes and affiliated files, and tags can be added to notes as well. Furthermore, it provides searching functionality, such as searching by title, date or tags.

## Architecture

This project is designed by DDD (Domain Driven Design) principles, which followed by this [article](https://programmingpercy.tech/blog/how-to-domain-driven-design-ddd-golang/).

```txt
Entity                      // Mutable Identifiable Structs
  |           ValueObject   // Immutable Unidentifiable Structs
  |                |
  └── Aggregate ───┘        // Combined set of Entities and Value objects, stored in Repositories
          |
      Repository            // A implementation of storing aggregates or other information
      (domain)
          |
      Factory               // A constructor to create complex objects and make creating new instance easier for the developers of other domains
          |
      Service               // A collection of repositories and sub-services that builds together the business flow
```

## Dependencies

- Gin: Web framework

- Gorm: ORM

- Mongo-driver: MongoDB driver

- Redis: Caching

- Zap + Lumberjack: Logging

- Swag: API Documentation
