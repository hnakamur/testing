#!/bin/bash

(./write-go | ./read.sh) &
pgrep -f 'read.sh' | xargs kill

