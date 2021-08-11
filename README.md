# ssh-separator
[![codecov](https://codecov.io/gh/mazrean/ssh-separator/branch/main/graph/badge.svg)](https://codecov.io/gh/mazrean/ssh-separator)
[![](https://github.com/mazrean/ssh-separator/workflows/CI/badge.svg)](https://github.com/mazrean/ssh-separator/actions)
[![](https://github.com/mazrean/ssh-separator/workflows/Release/badge.svg)](https://github.com/mazrean/ssh-separator/actions)
[![go report](https://goreportcard.com/badge/mazrean/ssh-separator)](https://goreportcard.com/report/mazrean/ssh-separator)

Tool to distribute ssh connections to containers for each user

## Requirement
* docker

## Usage
### Docker container
```
$ wget https://github.com/mazrean/ssh-separator/raw/main/docker-compose.yaml
$ docker compose up
```

## Environment Variables
|variable|description|example value|
|-|-|-|
|WELCOME|The string displayed upon successful login.|Login succeeded!|
|PROMETHEUS|If true, provide metrics for prometheus.|true|
|API_PORT|Port for REST API|3000|
|SSH_PORT|Port for ssh|2222|
|BADGER_DIR|Directory where user data is stored.|/var/lib/ssh-separator|
|IMAGE_NAME|Docker image for user container|mazrean/cpctf-ubuntu:latest|
|IMAGE_USER|Username in user containers.|ubuntu|
|IMAGE_CMD|Shell in user containers.|/bin/bash|
|CPU_LIMIT|The number of CPUs to allocate to user containers.|0.5|
|MEMORY_LIMIT|Memory limits for user containers.|1024|
|API_KEY|API key for REST API.|aeneexiene7uu3fie4pa|

## Licence
MIT
