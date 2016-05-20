#!/bin/bash

(./write.sh | ./read.sh) &
pgrep -f 'read.sh' | xargs kill

