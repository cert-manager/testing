#!/usr/bin/env bash

# +skip_license_check

# Run:
#    $ semversort 1.0 1.0-rc 1.0-patch 1.0-alpha
# or in GIT
#    $ semversort $(git tag)
# Using pipeline:
#    $ echo 1.0 1.0-rc 1.0-patch 1.0-alpha | semversort
#
# This script is from https://gist.githubusercontent.com/andkirby/54204328823febad9d34422427b1937b/raw/semversort.sh

set -o errexit
set -o pipefail
set -o nounset

if [ -t 0 ]; then
  versions_list=$@
else
  # catch pipeline output
  versions_list=$(cat)
fi

version_weight () {
  echo -e "$1" | tr ' ' "\n"  | sed -e 's:\+.*$::' | sed -e 's:^v::' | \
    sed -re 's:^[0-9]+(\.[0-9]+)+$:&-stable:' | \
    sed -re 's:([^A-Za-z])dev\.?([^A-Za-z]|$):\1.10.\2:g' | \
    sed -re 's:([^A-Za-z])(alpha|a)\.?([^A-Za-z]|$):\1.20.\3:g' | \
    sed -re 's:([^A-Za-z])(beta|b)\.?([^A-Za-z]|$):\1.30.\3:g' | \
    sed -re 's:([^A-Za-z])(rc|RC)\.?([^A-Za-z]|$)?:\1.40.\3:g' | \
    sed -re 's:([^A-Za-z])stable\.?([^A-Za-z]|$):\1.50.\2:g' | \
    sed -re 's:([^A-Za-z])pl\.?([^A-Za-z]|$):\1.60.\2:g' | \
    sed -re 's:([^A-Za-z])(patch|p)\.?([^A-Za-z]|$):\1.70.\3:g' | \
    sed -r 's:\.{2,}:.:' | \
    sed -r 's:\.$::' | \
    sed -r 's:-\.:.:'
}
tags_orig=(${versions_list})
tags_weight=( $(version_weight "${tags_orig[*]}") )

keys=$(for ix in ${!tags_weight[*]}; do
    printf "%s+%s\n" "${tags_weight[${ix}]}" ${ix}
done | sort -V | cut -d+ -f2)

for ix in ${keys}; do
    printf "%s\n" ${tags_orig[${ix}]}
done
