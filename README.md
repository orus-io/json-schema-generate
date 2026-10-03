# generate

Generates Go (golang) Structs and Validation code from JSON schema.

# Requirements

* [mise](https://mise.jdx.dev/) — manages the Go toolchain and project tasks

# Usage

Build the generator and generate test sources

```console
$ mise run generate
```

Run the test suite

```console
$ mise run test
```

Clean the generator binary and generated code

```console
$ mise run clean
```

Run the generator manually

```console
$ go build -o schema-generate ./cmd/schema-generate
$ ./schema-generate -p main -o output.go exampleschema.json
```

Flags

```
-o string              The output file for the schema.
-p string              The package that the structs are created in. (default "main")
-i string              A single file path (used for backwards compatibility).
-alwaysAcceptFalse     Any field will accept decoding 'false' and ignore it.
-useEmptyTypes         Use types with a empty types if non-required.
-schemaKeyRequired     Allow input files with no $schema key.
```

# Example

This schema

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "Example",
  "id": "http://example.com/exampleschema.json",
  "type": "object",
  "description": "An example JSON Schema",
  "properties": {
    "name": {
      "type": "string"
    },
    "address": {
      "$ref": "#/definitions/address"
    },
    "status": {
      "$ref": "#/definitions/status"
    }
  },
  "definitions": {
    "address": {
      "id": "address",
      "type": "object",
      "description": "Address",
      "properties": {
        "street": {
          "type": "string",
          "description": "Address 1",
          "maxLength": 40
        },
        "houseNumber": {
          "type": "integer",
          "description": "House Number"
        }
      }
    },
    "status": {
      "type": "object",
      "properties": {
        "favouritecat": {
          "enum": [
            "A",
            "B",
            "C"
          ],
          "type": "string",
          "description": "The favourite cat.",
          "maxLength": 1
        }
      }
    }
  }
}
```

generates

```go
package main

type Address struct {
  HouseNumber int `json:"houseNumber,omitempty"`
  Street string `json:"street,omitempty"`
}

type Example struct {
  Address *Address `json:"address,omitempty"`
  Name string `json:"name,omitempty"`
  Status *Status `json:"status,omitempty"`
}

type Status struct {
  Favouritecat string `json:"favouritecat,omitempty"`
}
```

See the [test/](./test/) directory for more examples.
