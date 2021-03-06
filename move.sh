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
    mv -v ~/oss/qt.go/qtqml/qgolistmodel.go ~/oss/qt.go/qtqml/_qgolistmodel.go
    # ~/oss/qt.go/ must be a soft link of $GOPATH/src/github.com/kitech/qt.go/
    mkdir -p ~/oss/qt.go/{qt5,qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}
    rm -f ~/oss/qt.go/{qt5,qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}/q*.go
    mkdir -p ~/oss/qt.go/{qtpositioning,qtwebchannel,qtwebenginecore,qtwebengine,qtwebenginewidgets}
    rm -f ~/oss/qt.go/{qtpositioning,qtwebchannel,qtwebenginecore,qtwebengine,qtwebenginewidgets}/q*.go
    mkdir -p ~/oss/qt.go/{qtsvg,qtmultimedia}
    rm -f ~/oss/qt.go/{qtsvg,qtmultimedia}/q*.go

    rm -f src/*/q*list.go
    # cp -a src/qtrt/*.go ~/oss/qt.go/qtrt/
    cp -a src/core/*.go ~/oss/qt.go/qtcore/
    cp -a src/gui/*.go ~/oss/qt.go/qtgui/
    cp -a src/widgets/*.go ~/oss/qt.go/qtwidgets/
    cp -a src/network/*.go ~/oss/qt.go/qtnetwork/
    cp -a src/qml/*.go ~/oss/qt.go/qtqml/
    mv -v ~/oss/qt.go/qtqml/_qgolistmodel.go ~/oss/qt.go/qtqml/qgolistmodel.go

    cp -a src/quick/*.go ~/oss/qt.go/qtquick/
    cp -a src/quickcontrols2/*.go ~/oss/qt.go/qtquickcontrols2/
    cp -a src/quickwidgets/*.go ~/oss/qt.go/qtquickwidgets/
    cp -a src/androidextras/*.go ~/oss/qt.go/qtandroidextras/
    cp -a src/winextras/*.go ~/oss/qt.go/qtwinextras/
    cp -a src/macextras/*.go ~/oss/qt.go/qtmacextras/

    # webengines
    cp -a src/positioning/*.go ~/oss/qt.go/qtpositioning/
    cp -a src/webchannel/*.go ~/oss/qt.go/qtwebchannel/
    cp -a src/webenginecore/*.go ~/oss/qt.go/qtwebenginecore/
    cp -a src/webengine/*.go ~/oss/qt.go/qtwebengine/
    cp -a src/webenginewidgets/*.go ~/oss/qt.go/qtwebenginewidgets/

    # multimedia
    cp -a src/multimedia/*.go ~/oss/qt.go/qtmultimedia/
    cp -a src/svg/*.go ~/oss/qt.go/qtsvg/
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
    mkdir -p ~/oss/qt.inline/src/androidextras
    mkdir -p ~/oss/qt.inline/src/{svg,multimedia}

    set +x

    mvbymd5 cxx src/core ~/oss/qt.inline/src/core
    mvbymd5 cxx src/gui ~/oss/qt.inline/src/gui
    mvbymd5 cxx src/widgets ~/oss/qt.inline/src/widgets
    mvbymd5 cxx src/network ~/oss/qt.inline/src/network
    mvbymd5 cxx src/qml ~/oss/qt.inline/src/qml
    mvbymd5 cxx src/quick ~/oss/qt.inline/src/quick
    mvbymd5 cxx src/quickcontrols2 ~/oss/qt.inline/src/quickcontrols2
    mvbymd5 cxx src/quickwidgets ~/oss/qt.inline/src/quickwidgets

    ### webengines
    mvbymd5 cxx src/positioning ~/oss/qt.inline/src/positioning
    mvbymd5 cxx src/webchannel ~/oss/qt.inline/src/webchannel
    mvbymd5 cxx src/webenginecore ~/oss/qt.inline/src/webenginecore
    mvbymd5 cxx src/webengine ~/oss/qt.inline/src/webengine
    mvbymd5 cxx src/webenginewidgets ~/oss/qt.inline/src/webenginewidgets

    ### extras
    mvbymd5 cxx src/androidextras ~/oss/qt.inline/src/androidextras

    ### multimedia
    mvbymd5 cxx src/svg ~/oss/qt.inline/src/svg
    mvbymd5 cxx src/multimedia ~/oss/qt.inline/src/multimedia

    #cp -a src/core/*.cxx ~/oss/qt.inline/src/core/
    #cp -a src/gui/*.cxx ~/oss/qt.inline/src/gui/
    #cp -a src/widgets/*.cxx ~/oss/qt.inline/src/widgets/
    # cp -a src/network/*.cxx ~/oss/qt.inline/src/network/
    # cp -a src/qml/*.cxx ~/oss/qt.inline/src/qml/
    # cp -a src/quick/*.cxx ~/oss/qt.inline/src/quick/

    # cp -a CMakeLists.txt ~/oss/qt.inline/
    # cp -a src/qihotfix.cpp ~/oss/qt.inline/src/
}

function mvrssrc()
{
    set +x

    mkdir -p ~/oss/qt.rs/src/{core,gui,widgets,network,qml,quick}
    # rm -f ~/oss/qt.rs/src/{core,gui,widgets,network,qml,quick}/q*.rs

    mvbymd5 rs src/core ~/oss/qt.rs/src/core
    mvbymd5 rs src/gui ~/oss/qt.rs/src/gui
    mvbymd5 rs src/widgets ~/oss/qt.rs/src/widgets
    # cp -a src/core/*.rs ~/oss/qt.rs/src/core/
    #cp -a src/gui/*.rs ~/oss/qt.rs/src/gui/
    #cp -a src/widgets/*.rs ~/oss/qt.rs/src/widgets/
    #cp -a src/network/*.rs ~/oss/qt.rs/src/network/
    #cp -a src/qml/*.rs ~/oss/qt.rs/src/qml/
    #cp -a src/quick/*.rs ~/oss/qt.rs/src/quick/

}

function mvcrsrc()
{
    set +x

    mkdir -p ~/oss/qt.cr/src/{qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}
    # rm -f ~/oss/qt.cr/src/{core,gui,widgets,network,qml,quick}/q*.cr

    mvbymd5 cr src/core ~/oss/qt.cr/src/qtcore
    # cp -a src/core/*.rs ~/oss/qt.cr/src/core/
    #cp -a src/gui/*.rs ~/oss/qt.cr/src/gui/
    #cp -a src/widgets/*.rs ~/oss/qt.cr/src/widgets/
    #cp -a src/network/*.rs ~/oss/qt.cr/src/network/
    #cp -a src/qml/*.rs ~/oss/qt.cr/src/qml/
    #cp -a src/quick/*.rs ~/oss/qt.cr/src/quick/

}

function mvdtsrc()
{
    set +x

    mkdir -p ~/oss/qtdart/src/{qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}
    # rm -f ~/oss/qtdart/src/{core,gui,widgets,network,qml,quick}/q*.dart

    mvbymd5 dart src/core ~/oss/qtdart/src/qtcore
    # cp -a src/core/*.rs ~/oss/qtdart/src/core/
    #cp -a src/gui/*.rs ~/oss/qtdart/src/gui/
    #cp -a src/widgets/*.rs ~/oss/qtdart/src/widgets/
    #cp -a src/network/*.rs ~/oss/qtdart/src/network/
    #cp -a src/qml/*.rs ~/oss/qtdart/src/qml/
    #cp -a src/quick/*.rs ~/oss/qtdart/src/quick/

}

function mvnimsrc()
{
    set +x

    mkdir -p ~/oss/qt5-nim/qt5/{qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}
    # rm -f ~/oss/qtdart/src/{core,gui,widgets,network,qml,quick}/q*.dart

    mvbymd5 nim src/core ~/oss/qt5-nim/qt5/qtcore
    mvbymd5 nim src/gui ~/oss/qt5-nim/qt5/qtgui
    # cp -a src/core/*.rs ~/oss/qtdart/src/core/
    #cp -a src/gui/*.rs ~/oss/qtdart/src/gui/
    #cp -a src/widgets/*.rs ~/oss/qtdart/src/widgets/
    #cp -a src/network/*.rs ~/oss/qtdart/src/network/
    #cp -a src/qml/*.rs ~/oss/qtdart/src/qml/
    #cp -a src/quick/*.rs ~/oss/qtdart/src/quick/

}

function mvvsrc()
{
    set +x

    mkdir -p ~/.vmodules/vqt/{qtcore,qtgui,qtwidgets,qtnetwork,qtqml,qtquick}
    #rm -f ~/.vmodules/vqt/qt{core,gui,widgets,network,qml,quick}/q*.v

    mvbymd5 v src/core ~/.vmodules/vqt/qtcore
    mvbymd5 v src/gui ~/.vmodules/vqt/qtgui
    mvbymd5 v src/widgets ~/.vmodules/vqt/qtwidgets
    #cp -a src/widgets/*.rs ~/oss/qt.cr/src/widgets/
    #cp -a src/network/*.rs ~/oss/qt.cr/src/network/
    #cp -a src/qml/*.rs ~/oss/qt.cr/src/qml/
    #cp -a src/quick/*.rs ~/oss/qt.cr/src/quick/

    rm -f ~/.vmodules/vqt/qt{core,gui,widgets}/*list.v
    rm -f ~/.vmodules/vqt/qt{core,gui,widgets}/*map.v
    cp -v src/gui/qpixmap.v ~/.vmodules/vqt/qtgui/
    rm -f ~/.vmodules/vqt/qt{core,gui,widgets}/*hash.v
    rm -f ~/.vmodules/vqt/qt{core,gui,widgets}/*set.v

    rm -f ~/.vmodules/vqt/qt{core,gui,widgets}/qfunctions.v

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
    crsrc)
        time mvcrsrc;
        ;;
    dtsrc)
        time mvdtsrc;
        ;;
    nimsrc)
        time mvnimsrc;
        ;;
    vsrc)
        time mvvsrc;
        ;;
    *)
        set +x
        echo "Unknown cmd: $cmd"
        help;
        ;;
esac


