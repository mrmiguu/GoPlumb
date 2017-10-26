package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mrmiguu/sock"
)

var (
	start, end int
)

func main() {
	F := sock.Wstring()
	R := sock.Rstring()
	E := sock.Rerror()

	code := js.Global.Get("document").Call("getElementById", "code")
	code.Set("onblur", func() {
		F <- js.Global.Get("window").Call("getSelection").String()
	})

	go func() {
		for {
			select {
			case res := <-R:
				swap(res)
			case err := <-E:
				swap("ERROR: " + err.Error())
			}
		}
	}()

	select {}
}

func swap(res string) {
	code := js.Global.Get("document").Call("getElementById", "code")
	text := code.Get("value").String()

	front := "\n"
	s, e := code.Get("selectionStart").Int(), code.Get("selectionEnd").Int()
	if s != e {
		front = ""
		start, end = s, e
	}

	code.Set("value", text[:start]+front+res+text[end:])

	start += len(front) + len(res)
	end = start
}
