
# Installation

*opmlsort* and *omplcat* are a command line programs run from a shell like Bash. You can find compiled
version in the [releases](https://github.com/rsdoiel/opml/releases/latest) 
in the Github repository in a zip file named like *opml-v0.0.3-release.zip*. Inside
the zip file look for the directory that matches your computer and copy that someplace
defined in your path (e.g. $HOME/bin). 

Compiled versions are available for Mac OS X (amd64 processor), Linux (amd64), Windows
(amd64) and Rapsberry Pi (both ARM6 and ARM7)

## Mac OS X

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green release zip file link and download
3. Open a finder window and find the downloaded file and unzip
4. Look in the unziped folder and find dist/macosx-amd64/opmlsort and dist/macosx-amd64/opmlcat
5. Drag (or copy) the *opmlsort* and *opmlcat* to a "bin" directory in your path
6. Open and "Terminal" and run `opmlsort -h` and `opmlcat -h`

## Windows

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green zip file link and download
3. Open the file manager find the downloaded file and unzip
4. Look in the unziped folder and find dist/windows-amd64/opmlsort.exe and dist/windows-amd64/opmlcat.exe
5. Drag (or copy) the *opmlsort.exe* *opmlcat.exe* to a "bin" directory in your path
6. Open Bash and and run `opmlsort -h` and `opmlcat -h`

## Linux

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green zip file link and download
3. find the downloaded zip file and unzip
4. In the unziped directory and find for dist/linux-amd64/opmlsort and dist/linux-amd64/opmlcat
5. copy the *opmlsort* and *opmlcat* to a "bin" directory (e.g. cp ~/Downloads/opml-binary-release/dist/linux-amd64/opml* ~/bin/)
6. From the shell prompt run `opmlsort -h` and `opmlcat -h`

## Raspberry Pi

If you are using a Raspberry Pi 2 or later use the ARM7 binary, ARM6 is only for the first generaiton Raspberry Pi.

1. Go to [github.com/rsdoiel/opml/releases/latest](https://github.com/rsdoiel/opml/releases/latest)
2. Click on the green zip file link and download
3. find the downloaded zip file and unzip
4. In the unziped directory and find for dist/raspberrypi-arm7/opmlsort and dist/raspberrypi-arm7/opmlcat
5. copy the *opmlsort* and *opmlcat* to a "bin" directory (e.g. cp ~/Downloads/opml-binary-release/dist/raspberrypi-arm7/opml* ~/bin/)
    + if you are using an original Raspberry Pi you should copy the ARM6 version instead
6. From the shell prompt run `opmlsort -h` and `opmlcat -h`
