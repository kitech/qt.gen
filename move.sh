#!/bin/sh

function help()
{
    echo "Usage:"
    echo "    move.sh <qil|gosrc|rssrc>"
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
    mkdir -p qt.go/src/{qt5,core,gui,widgets,network,qml,quick}
    rm -f qt.go/src/{qt5,core,gui,widgets,network,qml,quick}/q*.go

    # cp -a src/qtrt/*.go qt.go/src/qtrt/
    cp -a src/core/*.go qt.go/src/qt5/
    cp -a src/gui/*.go qt.go/src/qt5/
    cp -a src/widgets/*.go qt.go/src/qt5/
    cp -a src/network/*.go qt.go/src/qt5/
    cp -a src/qml/*.go qt.go/src/qt5/
    cp -a src/quick/*.go qt.go/src/qt5/
}

function mvqil()
{
    mkdir -p ~/oss/qt.inline/src/{qt5,core,gui,widgets,network,qml,quick}
    rm -f ~/oss/qt.inline/src/{qt5,core,gui,widgets,network,qml,quick}/q*.cxx

    cp -a src/core/*.{cxx,cmake} ~/oss/qt.inline/src/core/
    cp -a src/gui/*.{cxx,cmake} ~/oss/qt.inline/src/gui/
    cp -a src/widgets/*.{cxx,cmake} ~/oss/qt.inline/src/widgets/
    cp -a src/network/*.{cxx,cmake} ~/oss/qt.inline/src/network/
    cp -a src/qml/*.{cxx,cmake} ~/oss/qt.inline/src/qml/
    cp -a src/quick/*.{cxx,cmake} ~/oss/qt.inline/src/quick/

    # cp -a CMakeLists.txt ~/oss/qt.inline/
    cp -a src/qihotfix.cpp ~/oss/qt.inline/src/
}

function mvrssrc()
{
    mkdir -p ~/oss/qt.rs/src/{core,gui,widgets,network,qml,quick}
    rm -f ~/oss/qt.rs/src/{core,gui,widgets,network,qml,quick}/q*.rs

    cp -a src/core/*.rs ~/oss/qt.rs/src/core/
    cp -a src/gui/*.rs ~/oss/qt.rs/src/gui/
    cp -a src/widgets/*.rs ~/oss/qt.rs/src/widgets/
    cp -a src/network/*.rs ~/oss/qt.rs/src/network/
    cp -a src/qml/*.rs ~/oss/qt.rs/src/qml/
    cp -a src/quick/*.rs ~/oss/qt.rs/src/quick/
}

cmd=$1

set -x
case $cmd in
    qil)
        mvqil;
        ;;
    gosrc)
        mvgosrc;
        ;;
    rssrc)
        mvrssrc;
        ;;
    *)
        set +x
        echo "Unknown cmd: $cmd"
        help;
        ;;
esac


