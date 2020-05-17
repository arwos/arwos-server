#!/usr/bin/env bash

root=$(dirname $0)
dockercmd="docker-compose -f ${root}/docker-compose.yaml -p arwos-server"
