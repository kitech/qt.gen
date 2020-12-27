all:
	go build -v -i

run:
	go build -v -i
  # [QT_DIR=$HOME/Qt5.x.x/] ./qt.gen <c|go>
	./qt.gen

# for gov2, need remove .ast and then run two times
rerun:
	rm -fv ./qthdrsrc.ast

  # dont exist when error
	./qt.gen gov2 2>&1|grep "go:305" || true || true || true

	sleep 2
	./qt.gen gov2 2>&1|grep "go:305"

	sleep 3
	./move.sh gosrc

clean:
	rm -f qt.gen qthdrsrc.ast qthdrsrc.o
