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
    mkdir -p ~/oss/src/qt.go/{qt5,qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}
    rm -f ~/oss/src/qt.go/{qt5,qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}/q*.go

    # cp -a src/qtrt/*.go qt.go/src/qtrt/
    cp -a src/core/*.go ~/oss/src/qt.go/qtcore/
    cp -a src/gui/*.go ~/oss/src/qt.go/qtgui/
    cp -a src/widgets/*.go ~/oss/src/qt.go/qtwidgets/
    cp -a src/network/*.go ~/oss/src/qt.go/qtnetwork/
    cp -a src/qml/*.go ~/oss/src/qt.go/qtqml/
    cp -a src/quick/*.go ~/oss/src/qt.go/qtquick/
}

function mvqil()
{
    mkdir -p ~/oss/qt.inline/src/{qt5,core,gui,widgets,network,qml,quick}
    rm -f ~/oss/qt.inline/src/{qt5,core,gui,widgets,network,qml,quick}/q*.cxx

    cp -a src/core/*.cxx ~/oss/qt.inline/src/core/
    cp -a src/gui/*.cxx ~/oss/qt.inline/src/gui/
    cp -a src/widgets/*.cxx ~/oss/qt.inline/src/widgets/
    # cp -a src/network/*.cxx ~/oss/qt.inline/src/network/
    # cp -a src/qml/*.cxx ~/oss/qt.inline/src/qml/
    # cp -a src/quick/*.cxx ~/oss/qt.inline/src/quick/

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


