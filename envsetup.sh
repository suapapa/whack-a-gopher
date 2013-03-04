#!/bin/bash

export GOPATH=`realpath .`
export GOMAXPROC=`grep processor /proc/cpuinfo | awk '{field=$NF};END{print field+1}'`
