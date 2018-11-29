GUIDE TO CREATE GUI WITH GOTK3 AND GLADE
---------------------------------------

___________________

Installations gotk3
-------------------

```bash
$ go get github.com/gotk3/gotk3/...
```
### GNU / Linux

```bash
$ sudo apt-get install libgtk-3-dev
```

### Mac OS/X

```bash
$ brew install gtk-mac-integration gtk+3
```

### Windows

Reference: https://www.gtk.org/download/windows.php

- Download and install [MSYS2](http://www.msys2.org/)

- Install GTK+3 and its dependencies

```bash
$ pacman -S mingw-w64-x86_64-gtk3
```



Installation Glade
------------------

### GNU / Linux

```bash
$ sudo apt-get install glade
```

### Mac OS/X

```bash
$ brew install glade
```

### Windows

```bash
$ pacman -S mingw-w64-x86_64-glade
```

- Optionals

```bash
$ pacman -S mingw-w64-x86_64-devhelp
```

```bash
$ pacman -S mingw-w64-x86_64-toolchain base-devel
```

___________________
Example
-------

Build and run the go code from _ide_

```bash
$ go run ide_glade.go
```

___________________
Issues
-------

- `invalid flag in pkg-config --libs: -Wl,-luuid`

```bash
$ bash -c "sed -i -e 's/-Wl,-luuid/-luuid/g' C:/msys64/mingw64/lib/pkgconfig/gdk-3.0.pc"
```

- if `pkg-config` is not recognized as a command on Windows, remember to add `C:\msys64\mingw64\bin` to your PATH
