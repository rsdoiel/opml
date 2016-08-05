
# Installation

*opmlsort* is a command line program run from a shell like Bash. You can find compiled
version in the [releases](https://github.com/rsdoiel/opml/releases/latest) 
in the Github repository in a zip file named *opml-binary-release.zip*. Inside
the zip file look for the directory that matches your computer and copy that someplace
defined in your path (e.g. $HOME/bin). 

Compiled versions are available for Mac OS X (amd64 processor), Linux (amd64), Windows
(amd64) and Rapsberry Pi (both ARM6 and ARM7)

## Mac OS X

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green "opml-binary-release.zip" link and download
3. Open a finder window and find the downloaded file and unzip it (e.g. opml-binary-release.zip)
4. Look in the unziped folder and find dist/macosx-amd64/opmlsort
5. Drag (or copy) the *opmlsort* to a "bin" directory in your path
6. Open and "Terminal" and run `opmlsort -h`

## Windows

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green "opml-binary-release.zip" link and download
3. Open the file manager find the downloaded file and unzip it (e.g. opml-binary-release.zip)
4. Look in the unziped folder and find dist/windows-amd64/opmlsort.exe
5. Drag (or copy) the *opmlsort.exe* to a "bin" directory in your path
6. Open Bash and and run `opmlsort -h`

## Linux

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green "opml-binary-release.zip" link and download
3. find the downloaded zip file and unzip it (e.g. unzip ~/Downloads/opml-binary-release.zip)
4. In the unziped directory and find for dist/linux-amd64/opmlsort
5. copy the *opmlsort* to a "bin" directory (e.g. cp ~/Downloads/opml-binary-release/dist/linux-amd64/opmlsort ~/bin/)
6. From the shell prompt run `opmlsort -h`

## Raspberry Pi

If you are using a Raspberry Pi 2 or later use the ARM7 binary, ARM6 is only for the first generaiton Raspberry Pi.

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green "opml-binary-release.zip" link and download
3. find the downloaded zip file and unzip it (e.g. unzip ~/Downloads/opml-binary-release.zip)
4. In the unziped directory and find for dist/raspberrypi-arm7/opmlsort
5. copy the *opmlsort* to a "bin" directory (e.g. cp ~/Downloads/opml-binary-release/dist/raspberrypi-arm7/opmlsort ~/bin/)
    + if you are using an original Raspberry Pi you should copy the ARM6 version instead
6. From the shell prompt run `opmlsort -h`

