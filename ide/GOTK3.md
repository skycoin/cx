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

- Install Glade

```bash
$ pacman -S mingw-w64-x86_64-glade
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

http://ftp.gnome.org/pub/GNOME/binaries/win32/glade/

___________________
Example
-------

Build and run the go code from _ide_

```bash
$ go run ide_glade.go
```