package main

type Mangler interface {
	convTo(name string) string
	convFrom(mname string) string
}
