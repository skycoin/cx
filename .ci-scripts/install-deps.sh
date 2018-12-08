#!/bin/bash

if [[ $TRAVIS_OS_NAME == 'linux' ]]; then
    sudo apt-get update -qq
    sudo apt-get install --no-install-recommends -y \
                        libglib2.0-dev \
                        gtk+3.0 \
                        libgtk-3-dev \
                        libpango1.0-dev \
                        libcairo2-dev \
                        libperl-dev \
                        libgtk2.0-dev \
                        gir1.2-gtk-3.0 \
                        libxi-dev \
                        libgl1-mesa-dev \
                        libxrandr-dev \
                        libxcursor-dev \
                        libxinerama-dev \
                        xvfb
    export DISPLAY=:99.0
    sudo /usr/bin/Xvfb $DISPLAY 2>1 > /dev/null &
    export GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _| cut -d '_' -f 1-2)
    export Glib_VERSION=$(pkg-config --modversion glib-2.0)
    export Cairo_VERSION=$(pkg-config --modversion cairo)
    export Pango_VERSION=$(pkg-config --modversion pango)
    echo "GTK version ${GTK_VERSION} (Glib ${Glib_VERSION}, Cairo ${Cairo_VERSION}, Pango ${Pango_VERSION})"
fi
