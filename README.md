# gxsrvs/dtx types library

A Go library providing nullable type wrappers for database operations with JSON marshaling/unmarshaling support.

## Installation

```bash
go get github.com/gxsrvs/dtx
```

## Features

This library provides nullable variants of common Go types that are safe for database operations and JSON serialization:

- `NullString` - nullable string type
- `NullDate` - nullable date type with multiple format support
- `NullDateTime` - nullable datetime type  
- `NullTime` - nullable time type
- `NullInt16`, `NullInt32`, `NullInt64` - nullable integer types
- `NullFloat` - nullable float type
- `NullDecimal` - nullable decimal type using shopspring/decimal
- `NullBool` - nullable boolean type

## Usage

### Basic Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/gxsrvs/dtx/types"
)

func main() {
    // Create a nullable string with value
    ns := types.NewNullString("Hello, World!")
    fmt.Println("Value:", ns.ToString())
    fmt.Println("Is Empty:", ns.IsEmpty())
    
    // Create an empty nullable string
    emptyNs := types.NewNullStringEmpty()
    fmt.Println("Empty Value:", emptyNs.ToString())
    fmt.Println("Is Empty:", emptyNs.IsEmpty())
    
    // JSON marshaling
    jsonData, _ := json.Marshal(ns)
    fmt.Println("JSON:", string(jsonData))
    
    // JSON unmarshaling
    var unmarshaled types.NullString
    json.Unmarshal([]byte(`"Test Value"`), &unmarshaled)
    fmt.Println("Unmarshaled:", unmarshaled.ToString())
}
```

### Date Handling

```go
package main

import (
    "fmt"
    "time"
    "github.com/gxsrvs/dtx/types"
)

func main() {
    // Create nullable date
    now := time.Now()
    nd := types.NewNullDate(now)
    fmt.Println("Date:", nd.ToString())
    
    // Parse date from string (supports ISO and Russian formats)
    dateStr := "2023-12-25"
    parsed, err := types.ParseDateFromString(dateStr)
    if err == nil {
        nd2 := types.NewNullDate(*parsed)
        fmt.Println("Parsed Date:", nd2.ToString())
    }
}
```

### Database Integration

All types implement the `sql/driver.Valuer` and `sql.Scanner` interfaces for seamless database integration:

```go
package main

import (
    "database/sql"
    "github.com/gxsrvs/dtx/types"
)

type User struct {
    ID       int64             `db:"id"`
    Name     types.NullString  `db:"name"`
    Birthday types.NullDate    `db:"birthday"`
    Active   types.NullBool    `db:"active"`
}

func queryUser(db *sql.DB, id int64) (*User, error) {
    user := &User{}
    err := db.QueryRow("SELECT id, name, birthday, active FROM users WHERE id = ?", id).
        Scan(&user.ID, &user.Name, &user.Birthday, &user.Active)
    return user, err
}
```

## Supported Date Formats

The library supports multiple date formats for parsing:
- ISO format: `2006-01-02`
- Russian format: `02.01.2006`

## Dependencies

- `github.com/pkg/errors` - for error handling
- `github.com/shopspring/decimal` - for precise decimal operations

## Contributing

Contributions are welcome — please open an issue or a pull request.

**Project language: English.** All repository artefacts (source code, godoc
comments, Markdown documentation, commit messages, issue and pull request
descriptions, configuration files) must be written in English. This keeps the
project accessible to any Go developer who picks it up.

## License

This library is available under the MIT License.