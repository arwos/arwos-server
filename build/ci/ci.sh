#!/usr/bin/env bash

source $(dirname $0)/commands.sh

case $1 in
runapp)
  runapp
  ;;
alltests)
  alltests
  ;;
onetest)
  onetest
  ;;
*)
  echo "command not supported"
  ;;
esac
