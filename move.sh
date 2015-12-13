#!/bin/sh

clslst=$(cat src/core/widgets.rs |grep "mod q"|awk '{print $2}'|awk -F\; '{print $1}')

set -x
for cls in $clslst; do
    echo "src/core/$cls.rs"
    src="src/core/$cls.rs"
    dst="../../qt.rs/src/base"
    cp -v $src "$dst/"
done

