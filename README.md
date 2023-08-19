
<h1 align="center">bakery 🧁</h1>

[![Pipeline](https://github.com/kampanosg/bakery/actions/workflows/go.yml/badge.svg)](https://github.com/kampanosg/bakery/actions/workflows/go.yml)

Bakery is an (opinionated) alternative to the popular [GNU make](https://www.gnu.org/software/make/manual/make.html) tool. It uses YAML instead of a custom syntax and provides handy defaults, such as `help`, `version` and `metadata`. Its main goal is to make the developer experience more familiar and friendly while being fast and reliable.

## Installation
WIP...

## Usage

### The Bakefile


### Syntax
Below are the keywords that the `Bakefile` can contain

| keyword    | type                | optional | description                                                               |
| ---------- | ------------------- | -------- | ------------------------------------------------------------------------- |
| `version`  | `string`            | Y        | a user-defined version that can have any format e.g: `1` or `v1.2.3`      |
| `metadata` | `map[string]string` | Y        | any user-defined key-value pair. for example: `author: John`              |
| `defaults` | `[]string`          | Y        | a list of recipes that will be called if no recipe is passed at execution |
| `recipes`  | `[]Recipe`          | N        | the list of recipes, see the table below how to define them               |

The `Recipe` syntax is defined below
| keyword       | type       | optional | description                                                           |
| ------------- | ---------- | -------- | --------------------------------------------------------------------- |
| `description` | `string`   | Y        | a brief explanation of what the recipe does. used for the `bake help` |
| `steps`       | `[]string` | N        | the commands that the `recipe` executes                               |

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
  author: "George K <name@email.com>"
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
