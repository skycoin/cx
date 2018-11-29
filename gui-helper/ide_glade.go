package main

import (
    "log"
    "github.com/gotk3/gotk3/gtk"
)


func main() {
    gtk.Init(nil)

    b, err := gtk.BuilderNew()
    if err != nil {
        log.Fatal("Fatal: ", err)
    }

    err = b.AddFromFile("ide.glade")
    if err != nil {
        log.Fatal("Fatal: ", err)
    }

    obj, err := b.GetObject("window_main")
    if err != nil {
        log.Fatal("Fatal: ", err)
    }

    win := obj.(*gtk.Window)
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })

    win.ShowAll()

    gtk.Main()
}
