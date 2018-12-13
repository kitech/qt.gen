#!/bin/sh

# by latest archlinux qt4 package

INSDIR=$HOME

QTVER=4.8.7
rm -rf $INSDIR/Qt${QTVER}
mkdir -pv $INSDIR/Qt${QTVER}/${QTVER}/{gcc_64,Src}

cd $INSDIR/Qt${QTVER}/${QTVER}/gcc_64
pwd
ln -sv /usr/include/qt4 include

mkdir bin
ln -sv /usr/bin/qmake-qt4 bin/qmake
ln -sv /usr/bin/moc-qt4 bin/moc
ln -sv /usr/bin/rcc-qt4 bin/rcc

