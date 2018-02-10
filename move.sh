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
    cp -a src/quickcontrols2/*.go ~/oss/src/qt.go/qtquickcontrols2/
    cp -a src/quickwidgets/*.go ~/oss/src/qt.go/qtquickwidgets/
    cp -a src/androidextras/*.go ~/oss/src/qt.go/qtandroidextras/
    cp -a src/winextras/*.go ~/oss/src/qt.go/qtwinextras/
    cp -a src/macextras/*.go ~/oss/src/qt.go/qtmacextras/
}

function mvbymd5()
{
    extstr=$1
    srcdir=$2
    dstdir=$3
    if [ -z "$srcdir" ] || [ -z "$dstdir" ]; then
        echo "$srcdir => $dstdir"
        echo "empty dir"
        exit
    fi

    files=$(ls $srcdir/*.$extstr)
    for file in $files ; do
        # echo "123, $file"
        bname=$(basename $file)
        dfpath="$dstdir/$bname"
        needcp="yes"
        if [ -f "$dfpath" ]; then
            srcmd5=$(md5sum $file|awk '{print $1}')
            dstmd5=$(md5sum $dfpath|awk '{print $1}')
            if [ "$srcmd5" == "$dstmd5" ]; then
                needcp="no"
            fi
        fi
        if [ "$needcp" == "yes" ]; then
            true
            echo "install -m 0644 $file $dstdir/$bname"
            install -m 0644 "$file" "$dstdir/$bname"
        fi
    done
}

function mvqil()
{
    mkdir -p ~/oss/qt.inline/src/{qt5,core,gui,widgets,network,qml,quick,quickcontrols2,quickwidgets}
    # rm -f ~/oss/qt.inline/src/{qt5,core,gui,widgets,network,qml,quick}/q*.cxx

    set +x
    mvbymd5 cxx src/core ~/oss/qt.inline/src/core
    mvbymd5 cxx src/gui ~/oss/qt.inline/src/gui
    mvbymd5 cxx src/widgets ~/oss/qt.inline/src/widgets
    mvbymd5 cxx src/network ~/oss/qt.inline/src/network
    mvbymd5 cxx src/qml ~/oss/qt.inline/src/qml
    mvbymd5 cxx src/quick ~/oss/qt.inline/src/quick
    mvbymd5 cxx src/quickcontrols2 ~/oss/qt.inline/src/quickcontrols2
    mvbymd5 cxx src/quickwidgets ~/oss/qt.inline/src/quickwidgets

    #cp -a src/core/*.cxx ~/oss/qt.inline/src/core/
    #cp -a src/gui/*.cxx ~/oss/qt.inline/src/gui/
    #cp -a src/widgets/*.cxx ~/oss/qt.inline/src/widgets/
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
        time mvqil;
        ;;
    gosrc)
        time mvgosrc;
        ;;
    rssrc)
        time mvrssrc;
        ;;
    *)
        set +x
        echo "Unknown cmd: $cmd"
        help;
        ;;
esac


