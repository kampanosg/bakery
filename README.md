
# bakery 🧁 

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kampanosg/bakery/go.yml?style=for-the-badge&logo=go)

Bakery is an (opinionated) automation tool designed to simplify repetitive, time-consuming and error prone commands. It uses YAML instead of a propriatery language and provides handy defaults, such as `help` and `version`. Its main goal is to make the developer experience more familiar, flexible and friendly while being fast and reliable.

https://github.com/kampanosg/bakery/assets/30287348/deb1b130-042f-4e18-b792-2103fd5c008f

## Installation
Bakery can be installed with the following methods:
* [Go](#go)

### Go
Use the `go` CLI tool to install Bakery:
```bash
go install github.com/kampanosg/bakery@latest
```

## Usage

### The Bakefile
Bakery, by default, looks for a `Bakefile` from the directory that is executed. The `Bakefile`, should follow the YAML structure and contain the keywords and syntax described bellow. 

The `--file` flag can be passed to use a `Bakefile` in a different location. For example, `bake --file /path/to/Bakefile version`

### Syntax
#### Keywords
Below are the keywords that the `Bakefile` can contain

| keyword    | type                | optional | description                                                               |
| ---------- | ------------------- | -------- | ------------------------------------------------------------------------- |
| `version`  | `string`            | Y        | a user-defined version that can have any format e.g: `1` or `v1.2.3`      |
| `metadata` | `map[string]string` | Y        | any user-defined key-value pair. for example: `author: John`              |
| `defaults` | `[]string`          | Y        | a list of recipes that will be called if no recipe is passed at execution |
| `recipes`  | `[]Recipe`          | N        | the list of recipes, see the table below how to define them               |

#### Recipe
The `Recipe` syntax is defined below

| keyword       | type       | optional | description                                                           |
| ------------- | ---------- | -------- | --------------------------------------------------------------------- |
| `description` | `string`   | Y        | a brief explanation of what the recipe does. used for the `bake help` |
| `steps`       | `[]string` | N        | the commands that the `recipe` executes                               |

ℹ️ A `step` can reference a `recipe`

### Builtins
The following builtin functions are available with `bake`

| command   | description                                                                |
| --------- | -------------------------------------------------------------------------- |
| `version` | prints out the version of the `Bakefile` if it has been defined            |
| `help`    | prints out the available recipes in the `Bakefile` and their `description` |


### Example Bakefile
The following is an example `Bakefile` for a Go project. It includes examples for all the available syntax.

```yaml
version: v1.2.3
metadata:
  author: "Darth Vader <vader@empire.org>"
defaults:
  - build
recipes:
  build:
    description: "builds the project"
    steps:
      - "rm app"
      - "go build -p app ./..."
  test:
    description: "tests the project"
    steps:
      - "go test ./..."
  run:
    description: "builds and tests the project"
    steps:
      - test
      - build
      - "./app"
```
