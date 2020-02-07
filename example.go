// +build run

package main

import (
	"os"

	"github.com/zetamatta/go-windows-consolefontpixel"
)

func main1() error {
	dc, err := consolefontpixel.OpenDC()
	if err != nil {
		return err
	}
	defer dc.Close()

	for _, s := range os.Args[1:] {
		w, h, err := dc.GetTextExtentPoint(s)
		if err != nil {
			return err
		}
		println(s, ":width=", w, ",height=", h)
	}
	w, h, err := consolefontpixel.GetCurrentConsoleFont()
	if err != nil {
		return err
	}
	println("font width=", w, " height=", h)
	return nil
}

func main() {
	if err := main1(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
