package main

import "testing"

func TestQueryComment(t *testing.T) {
	queryCommentFromFile("/home/me/Qt5.10.1/Docs/Qt-5.10.1/qtgui/qwindow.html", "setMinimumSize")
}
