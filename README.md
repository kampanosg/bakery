# bakery 🧁 

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kampanosg/bakery/go.yml?style=for-the-badge&logo=go)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kampanosg/bakery/sec.yml?style=for-the-badge&logo=go&label=security)

Bakery is an (opinionated) automation tool designed to simplify repetitive, time-consuming and error prone commands. It uses YAML instead of a propriatery language and provides handy defaults, such as `help` and `version`. Its main goal is to make the developer experience more familiar, flexible and friendly while being fast and reliable.

https://github.com/kampanosg/bakery/assets/30287348/379cbb67-13a5-473e-992c-52209b0e2042

## Installation
Bakery can be installed with the following methods:
* [Go](#go)

### Go
Use the `go` CLI tool to install Bakery:
```bash
go install github.com/kampanosg/bakery/cmd/bake@latest
```

## Usage

### CLI
To run a `Bakefile` and its available recipes, use the `bake` CLI tool. For example:
```bash
bake version
bake build
```

#### Flags
The CLI ships with a set of command line flags that can be used to extract information from the tool or modify its behaviour.

| flag       | description                                                                     | default |
| ---------- | ------------------------------------------------------------------------------- | ------- |
| `-version` | prints the current version of the tool                                          |         |
| `-verbose` | controls the verbosity. if set to `true` the tool prints a lot more to `stdout` | `false` |
| `-help`    | prints a help message with all the available flags                              |         |

To use the flags, you can pass them directly to the CLI:

```bash
bake -version
bake -verbose build
bake -help
```

> :info: The `-verbose` flag can be used with recipes

### The Bakefile
Bakery, by default, looks for a `Bakefile` from the directory that is executed. The `Bakefile`, should follow the YAML structure and contain the keywords and syntax described bellow. 

The `--file` flag can be passed to use a `Bakefile` in a different location. For example, `bake --file /path/to/Bakefile version`

### Syntax
#### Keywords
Below are the keywords that the `Bakefile` can contain

| keyword    | type                | optional | description                                                                    |
| ---------- | ------------------- | -------- | ------------------------------------------------------------------------------ |
| `version`  | `string`            | Y        | a user-defined version that can have any format e.g: `1` or `v1.2.3`           |
| `metadata` | `map[string]string` | Y        | any user-defined key-value pair. for example: `author: John`                   |
| `variables`| `map[string]string` | Y        | user defined variables that can be used in the recipes with the `:var:` syntax |
| `defaults` | `[]string`          | Y        | a list of recipes that will be called if no recipe is passed at execution      |
| `recipes`  | `[]Recipe`          | N        | the list of recipes, see the table below how to define them                    |
| `help`     | `string`            | Y        | a custom message that overrides the built-in `bake help`                       |

#### Recipe
The `Recipe` syntax is defined below

| keyword       | type       | optional | description                                                           |
| ------------- | ---------- | -------- | --------------------------------------------------------------------- |
| `description` | `string`   | Y        | a brief explanation of what the recipe does. used for the `bake help` |
| `steps`       | `[]string` | N        | the commands that the `recipe` executes                               |
| `private`     | `bool`     | Y        | if set to true it cannot be called directly from the command line     |

ℹ️ A `step` can reference a `recipe`

#### Special Characters
Some special keywords that are available to the `Bakefile`

| character | usage | description                                                              |
| --------- | ----- | ------------------------------------------------------------------------ |
| `^`       | steps | prefix a step with this keyword to ignore errors and continue the recipe |

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
variables:
  name: app
defaults:
  - build
recipes:
  build:
    description: "builds the project"
    steps:
      - "rm :name:"
      - "go build -p :name: ./..."
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
  ignore-fail:
    private: true
    description: "prefix a step with ^ to ignore the failure"
    steps:
        - "^lss -abcdf ./not/valid/dir"
        - "echo 'i will execute fine'"
help: |
  This is an example help message. It will override the built-in `bake help`

  And it can be multiline as well!
```
