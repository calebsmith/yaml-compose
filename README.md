# yaml-compose

A command-line tool for composing several YAML files together.

## Rationale

Many tools such as k8s, helm, and docker-compose use YAML configuration files. Some tools such as Azure Pipelines have templating capability but
it is otherwise difficult to keep these files DRY and organized.

Namely, it is often desirable to:
1. Loading - Load a YAML file as data and use these as template variables
2. Injection - Directly include the contents of one YAML file in another.

yaml-compose has simple mechanisms for each of these and they can work recursively and in tandem. i.e. One can inject a template that injects another, loads yet another and so on.

(See "Usage" below for a practical example for docker-compose files)

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