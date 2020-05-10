#!/usr/bin/env bash

source `dirname $0`/vars.sh

alltests(){
  ${dockercmd} exec app go test ./...
}

onetest(){
  ${dockercmd} exec app go test -run "$2" ./...
}

runapp(){
  ${dockercmd} down
  ${dockercmd} up
  ${dockercmd} down
}