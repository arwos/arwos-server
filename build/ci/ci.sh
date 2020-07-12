#!/usr/bin/env bash

source $(dirname $0)/commands.sh

case $1 in
    envup)
        envup;;
    envdown)
        envdown;;
    alltests)
        alltests;;
    onetest)
        onetest;;
    runapp)
        runapp;;
    *)
        echo "command not supported";;
esac
