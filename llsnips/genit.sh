
QTINC="-I/usr/include/qt/ -I/usr/include/qt/QtCore/"
set -x
clang++ -S -emit-llvm -o snip1.x64.ll  $QTINC snip1.cpp
clang++ -S -emit-llvm -m32 -o snip1.x32.ll  $QTINC snip1.cpp

# grep foo

