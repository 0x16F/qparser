## Overview

The `qparser` module is a powerful query parser and builder for Go, designed to work seamlessly with the GORM library and web frameworks like Fiber. It simplifies the process of parsing query parameters from web requests and applying them to database queries, allowing for flexible and dynamic data retrieval based on client requests.

## Features

- **Easy Query Parsing**: Automatically parse query parameters into GORM-compatible queries.
- **Flexible Query Options**: Supports various operators like equals, not equal, greater than, less than, like, and range queries.
- **Seamless Integration**: Designed to work effortlessly with GORM and web frameworks like Fiber.

## Installation

To install `qparser`, use the following go get command:

```
go get github.com/0x16F/qparser
```

## Usage

### Defining Request Structs

Define your request struct with tags specifying how to parse each field. Use the `query` tag to indicate which query parameter it corresponds to.

```go
type Request struct {
	Name  string `query:"name"`
	Email string `query:"email"`
}
```

### Parsing and Applying Queries

Within your request handler, parse the request into a struct, then use `qparser` to generate query options and apply them to your database queries.

Example with Fiber and GORM:

```go
package main

import (
	"github.com/0x16F/qparser"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func main() {
	db, err := gorm.Open(&gorm.Config{})
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		var req Request

		if err := c.QueryParser(&req); err != nil {
			return err
		}

		options, err := qparser.ParseStruct(&req)
		if err != nil {
			return err
		}

		var users []User

		if err := options.Apply(db.WithContext(c.Context()).Model(&User{})).Find(&users).Error; err != nil {
			return err
		}

		return c.JSON(users)
	})

	app.Listen(":3000")
}
```

## Supported Operators

`qparser` supports a variety of operators for query building:

- `eq`: Equals
- `neq`: Not equals
- `gt`: Greater than
- `gte`: Greater than or equal to
- `lt`: Less than
- `lte`: Less than or equal to
- `like`: Like (for pattern matching)
- `rng`: Range (for between queries)

Each operator is mapped to its SQL equivalent, ensuring accurate query construction.

### Examples of URL Query Parameters and Their SQL Representations

#### Equals (`eq`)

**HTTP Request:**

```
example.com/users?status=eq:1
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE status = 1;
```

#### Not Equals (`neq`)

**HTTP Request:**

```
example.com/users?status=neq:1
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE status <> 1;
```

#### Greater Than (`gt`)

**HTTP Request:**

```
example.com/users?age=gt:30
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE age > 30;
```

#### Greater Than or Equal To (`gte`)

**HTTP Request:**

```
example.com/users?age=gte:30
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE age >= 30;
```

#### Less Than (`lt`)

**HTTP Request:**

```
example.com/users?age=lt:30
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE age < 30;
```

#### Less Than or Equal To (`lte`)

**HTTP Request:**

```
example.com/users?age=lte:30
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE age <= 30;
```

#### Like (`like`)

**HTTP Request:**

```
example.com/users?name=like:John
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE name ILIKE '%John%';
```

The `%` symbols are added by the application to conduct a pattern match.

#### Range (`rng`)

**HTTP Request:**

```
example.com/users?createdAt=rng:2020-01-01 to 2020-12-31
```

**SQL Representation:**

```sql
SELECT * FROM users WHERE createdAt BETWEEN '2020-01-01' AND '2020-12-31';
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.