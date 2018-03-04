#!/bin/sh

set -x

SRCDIR=/home/me/tmp/vqtenv/qt-everywhere-src-5.10.1
SYNCQT=$SRCDIR/qtbase/bin/syncqt.pl
qtver=5.10.1
DSTDIR=./qtheaders/

topdirs=$(ls -l $SRCDIR | grep "^d" | grep -v coin | grep -v gnuwin32 | grep -v qtdoc | awk '{print $9}')

mkdir $DSTDIR -pv
for topdir in $topdirs; do
    echo "$topdir"
    if [ -d $SRCDIR/$topdir/include ]; then
        ls "$SRCDIR/$topdir/include/"
        # cp -a $SRCDIR/$topdir/include/* $DSTDIR/
        echo "$SYNCQT -copy -windows -showonly -private -version $qtver -outdir $DSTDIR/ $SRCDIR/$topdir"
        $SYNCQT -copy -windows -private -version $qtver -outdir $DSTDIR/ $SRCDIR/$topdir
        subdirs=$(ls "$SRCDIR/$topdir/include/")
        for subdir in $subdirs; do
            echo $subdir
            echo "touch $DSTDIR/include/$subdir/${subdir}Depends"
            touch $DSTDIR/include/$subdir/${subdir}Depends
            touch $DSTDIR/include/$subdir/${subdir,,}-config.h
        done
    fi
done

cp -v $SRCDIR/qtbase/src/corelib/global/qglobal.h $DSTDIR/include/QtCore/
cp -v $SRCDIR/qtbase/src/corelib/global/qconfig-bootstrapped.h $DSTDIR/include/QtCore/
touch $DSTDIR/include/QtCore/qconfig.h
mkdir -p $DSTDIR/include/CoreFoundation/
touch $DSTDIR/include/CoreFoundation/CoreFoundation
touch $DSTDIR/include/jni.h


