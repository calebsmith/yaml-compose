# yaml-compose

A command-line tool for composing several YAML files together.

Warning: This project is in an initial proof-of-concept phase and should not be used for production. Use at your own risk.

(I am open to feedback on design ideas or feature requests)

## Rationale

Many tools such as k8s, helm, and docker-compose use YAML configuration files extensively. Some tools such as Azure Pipelines have templating capability but it is otherwise difficult to keep these files DRY and organized and not available for all such tools.

Namely, it is often desirable to:
1. Load a YAML file as data and use these as template variables
2. Directly include the contents of one YAML file in another.

In yaml-compose, these techniques are called "loading" and "injecting" respectively.

yaml-compose has simple mechanisms for each of these and they can work recursively and in tandem. i.e. One can inject a template that injects another, loads yet another and so on.

(See "Usage" below for more details and examples)

## Installation

Install from source:
```
git clone git@github.com:calebsmith/yaml-compose
cd yaml-compose
make build
```
From here, copy `yaml-compose` to your PATH, or otherwise append the yaml-compose directory to your PATH.

TODO: Add release automation and installation instructions from Github releases

## Usage

The syntax for yaml-compose is as follows:
1. `{# vars.yaml #}` - loads a file named "vars.yaml"
2. `{$ filename.yaml $}` - injects the contents of "filename.yaml"

Some further considerations and examples:
* Loading a file *only* loads the data therein. The directive is ellided from the file in which it was found. This data can then be used with `{{.VarName}}` per normal Golang text/template syntax.
* Injecting a file maintains any prefix string on the line with the directive. Any contents that are loaded in maintain the level of indentation of that call site.

## Examples

### Variable Loading

```
# vars.yaml
Host: localhost
Port: 3000
```

```
# main.yaml

server:
  host: {{.Host}}
  port: {{.Port}}
```

Results in:
```
server:
  host: localhost
  port: 3000
```

### File Injection

```
# other.yaml
key3: value3
key4: value4
```
```
# main.yaml
data:
  - key1: value1
    key2: value2
  - {$ other.yaml $}
```

Will result in:
```
data:
  - key1: value1
    key2: value2
  - key3: value3
    key4: value4
```

### Realistic Example

See `examples/docker-compose` to see an example that creates a docker-compose file. Running `yaml-compose infra.yaml > result.yaml` in that folder will produce the `result.yaml` file to be used with docker-compose. This allows for separation of concerns and docker-compose files can be "composed" together with yaml-compose in this way.

```
# examples/docker-compose/infra.yaml

{# config.yaml #}
version: '3'
services:
  {$ partials/postgres.yaml $}
  {$ partials/redis.yaml $}
```

N.B. - One may be tempted to use process substitution such as the following:
```
docker-compose -f <(yaml-compose infra.yaml) up -d
```

However, in practice, docker-compose leverages the working directory of the compose file for certain needs such as volume mounting.

As a workaround, consider a shell function or alias such as:

```
# ~/.zprofile

# compose docker-compose and yaml-compose together
function dcc() {
  file="$1"
  shift
  TMPPREFIX=$PWD docker-compose -f =(yaml-compose "$file") "$@"
};
```

Then, the following works as intended:
`dcc infra.yaml up -d`

## Testing

To test yaml-compose use:

```
make test
```

This uses the `cram` Python library for functional testing of command line programs. To install it, use:

```
make build_test
```

The test make target may eventually include Go unit tests as well.