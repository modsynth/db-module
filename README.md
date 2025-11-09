# DB Module

> Database abstraction layer with Repository pattern using GORM

Part of the [Modsynth](https://github.com/modsynth) ecosystem.

## Features

- Multi-database support (PostgreSQL, MySQL, SQLite)
- Generic Repository pattern with type safety
- Connection pooling and management
- Transaction support
- Health checks and statistics
- Auto-migration support

## Installation

```bash
go get github.com/modsynth/db-module
```

## Quick Start

### Database Connection

```go
package main

import (
    "github.com/modsynth/db-module"
    "gorm.io/driver/postgres"
)

func main() {
    config := &db.Config{
        Driver:          "postgres",
        DSN:             "host=localhost user=myuser password=mypass dbname=mydb port=5432",
        MaxOpenConns:    100,
        MaxIdleConns:    10,
    }

    database, err := db.New(config, postgres.Open(config.DSN))
    if err != nil {
        panic(err)
    }
    defer database.Close()
}
```

### Repository Pattern

```go
package main

import (
    "context"
    "github.com/modsynth/db-module/repository"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Email string
}

func main() {
    // Create repository
    userRepo := repository.New[User](database.DB)

    // Create
    user := &User{Name: "John", Email: "john@example.com"}
    userRepo.Create(context.Background(), user)

    // Find by ID
    var found User
    userRepo.FindByID(context.Background(), 1, &found)

    // Find all
    users, _ := userRepo.FindAll(context.Background())

    // Update
    found.Name = "Jane"
    userRepo.Update(context.Background(), &found)

    // Delete
    userRepo.Delete(context.Background(), &found)
}
```

## Supported Databases

- PostgreSQL - `gorm.io/driver/postgres`
- MySQL - `gorm.io/driver/mysql`
- SQLite - `gorm.io/driver/sqlite`

## Version

Current version: `v0.1.0`

## License

MIT
