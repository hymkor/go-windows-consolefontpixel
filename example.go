// +build run

package main

import (
	"os"

	"github.com/zetamatta/go-windows-consolefontpixel"
)

func main() {
	for _, s := range os.Args[1:] {
		w, h, err := consolefontpixel.GetPixelSize(s)
		if err != nil {
			println(err.Error())
			return
		}
		println(s, ":width=", w, ",height=", h)
	}
	w, h, err := consolefontpixel.GetFontSize()
	if err != nil {
		println(err.Error())
		return
	}
	println("font width=", w, " height=", h)
}
