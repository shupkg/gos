#!/bin/sh
cd $(dirname $0)
find _templates -type f -exec wtool embed -i github.com/shupkg/gos/embed -f -o . {} \;
