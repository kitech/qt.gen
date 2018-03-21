#!/bin/sh

# only WRITABLE property
hdrfiles=$(ls /usr/include/qt/QtWidgets/q*.h | grep -v "_p.h" | grep -v "\-config.h")

for f in $hdrfiles; do
    echo "$f"
    # cat $f | grep "Q_PROPERTY(" | grep "DESIGNABLE"
    cat $f | grep "Q_PROPERTY(" | grep "WRITE"
done
