#!/bin/sh

# clslst=$(cat src/core/widgets.rs |grep "mod q"|awk '{print $2}'|awk -F\; '{print $1}')

# set -x
# for cls in $clslst; do
#    echo "src/core/$cls.rs"
#    src="src/core/$cls.rs"
#    dst="../../qt.rs/src/base"
#    cp -v $src "$dst/"
# done

function mvgosrc()
{
    mkdir -p qt.go/src/{qt5,core,gui,widgets}
    cp -a src/core/*.go qt.go/src/qt5/
    cp -a src/gui/*.go qt.go/src/qt5/
    cp -a src/widgets/*.go qt.go/src/qt5/
}

mvgosrc;
