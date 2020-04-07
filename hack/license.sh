#!/usr/bin/env bash

set -exo pipefail

export LC_ALL=C

cd $(dirname $0)/..
ROOT=$(git rev-parse --show-toplevel)
LICENSE_FILE=LICENSE-3rdparty.csv
FILE=${1:-$LICENSE_FILE} 

echo Component,Origin,License > $FILE
echo 'core,"github.com/frapposelli/wwhrd",MIT' >> $FILE
unset grep
$ROOT/bin/wwhrd list --no-color |& grep "Found License" | awk '{print $6,$5}' | sed -E "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]//g" | sed s/" license="/,/ | sed s/package=/core,/ | sort >> $FILE
