all:
	go build -v -i

run:
	go build -v -i
  # [QT_DIR=$HOME/Qt5.x.x/] ./qt.gen <c|go>
	./qt.gen

clean:
	rm -f qt.gen qthdrsrc.ast qthdrsrc.o
