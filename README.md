# riff

<!-- [![Build Status](https://travis-ci.org/gimke/riff.svg?branch=master)](https://travis-ci.org/gimke/riff) [![GoDoc](https://godoc.org/github.com/teatak/riff?status.svg)](https://godoc.org/github.com/teatak/riff) [![Join the chat at https://gitter.im/gimke/riff](https://badges.gitter.im/gimke/riff.svg)](https://gitter.im/gimke/riff?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) -->

Riff is a tool for service discovery and configuration.

* Service Discovery
* Health Checking
* Web console

## Quick Start
If you wish to work on Riff, you'll first need Go installed (version 1.9+ is required).

If you need Web Console, you'll need npm installed dependencies            

Next clone this repository and build Riff or download from [release](https://github.com/teatak/riff/releases)
```bash
make
./riff -v
```
## Commands(CLI)
```bash
$ ./riff

Usage: riff [--version] <command> [<args>]

Available commands are:

  daem        Run Riff as service 
  query       Query
  quit        Quit Riff
  reload      Reload Riff config
  restart     Restart service
  run         Run Riff
  start       Start service
  stop        Stop service
  update      Update Riff
  version     Prints the Riff version


```

### run
run Riff
```bash
$ ./riff run -h

Usage: run [options]

  Run riff service

Options:

  -name       Node name
  -dc         DataCenter name
  -http       Http address of riff (-http 127.0.0.1:8610)
  -rpc        RPC address of riff (-rpc [::]:8630)
  -join       Join RPC address (-join 192.168.1.1:8630,192.168.1.2:8630,192.168.1.3:8630)

```
## API

### graphql 
> POST /api

you can use explorer in Web console to discover api

```graphql
{
  service(name: "mongod") {
    name
    nodes {
      name
    }
  }
}

```
```bash
curl --request POST \
  --url http://localhost:8610/api \
  --header 'content-type: application/json' \
  --data '{"query":"{\n  service(name: \"mongod\") {\n    name\n    nodes {\n      name\n    }\n  }\n}\n"}'
```
### logs

> GET /api/logs

watch logs
```bash
curl --request GET \
  --url http://localhost:8610/api/logs
```
### watch

> POST /api/watch?name={watchName}&type={node|service}

watch node or service

```bash
curl --request POST \
  --url 'http://localhost:8610/api/watch?type=service&name=mongod' \
  --header 'content-type: application/json' \
  --data '{"query":"{\n  service(name: \"mongod\") {\n    name\n    nodes {\n      name\n    }\n  }\n}\n"}'
```

## Web Console

> http://localhost:8610

![Riff console](https://raw.githubusercontent.com/gimke/riff/gh-pages/images/screen.png)

## Service Config
config files in config/*.yml

ping.yml config file
```yaml
#name: service name
#port: service port
#env:
#  - CART_MODE=release

#command:
#  - ./home/cartdemo/cartdemo

#pid_file: ./home/cartdemo/cartdemo.pid
#std_out_file: ./home/cartdemo/logs/out.log
#std_err_file: ./home/cartdemo/logs/err.log
#grace: true
#run_at_load: false
#keep_alive: false

#deploy:
#  service_path: (service path e.g /home/riff)
#  provider: github (only support github gitlab)
#  token: Personal access tokens (visit https://github.com/settings/tokens or https://gitlab.com/profile/personal_access_tokens and generate a new token)
#  repository: repository address (https://github.com/teatak/cartdemo)
#  version: branchName (e.g branch:master), latest release (e.g release:latest or tag:latest）,or (content:master:version.txt)
#  payload: payload url when update success

name: ping
env:
  - CART_MODE=release

command:
  - ping
  - 192.168.3.1

grace: false
run_at_load: true

```
