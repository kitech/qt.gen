#!/bin/sh

function help()
{
    echo "Usage:"
    echo "    move.sh <qi|gosrc>"
}

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
    # cp -a src/qtrt/*.go qt.go/src/qtrt/
    cp -a src/core/*.go qt.go/src/qt5/
    cp -a src/gui/*.go qt.go/src/qt5/
    cp -a src/widgets/*.go qt.go/src/qt5/
}

function mvqi()
{
    cp -a src/core/*.{cxx,cmake} ~/oss/qt.inline/src/core/
    cp -a src/gui/*.{cxx,cmake} ~/oss/qt.inline/src/gui/
    cp -a src/widgets/*.{cxx,cmake} ~/oss/qt.inline/src/widgets/

    # cp -a CMakeLists.txt ~/oss/qt.inline/
    cp -a src/qihotfix.cpp ~/oss/qt.inline/src/
}

cmd=$1

set -x
case $cmd in
    qi)
        mvqi;
        ;;
    gosrc)
        mvgosrc;
        ;;
    *)
        set +x
        echo "Unknown cmd: $cmd"
        help;
        ;;
esac


