Installation for development of **opml**
===========================================

**opml** A library for working with OPML XML files and example command line programs.

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the following command to install latest table release. This may work for macOS, Linux and if you’re using Windows with the Unix subsystem. This would be run from your shell (e.g. Terminal on macOS).

~~~shell
curl https://rsdoiel.github.io/opml/installer.sh | sh
~~~

This will install the programs included in opml in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command below.

~~~ps1
irm https://rsdoiel.github.io/opml/installer.ps1 | iex
~~~

### If your are running macOS or Windows

You may get security warnings if you are using macOS or Windows. See the notes for the specific operating system you're using to fix issues.

- [INSTALL_NOTES_macOS.md](INSTALL_NOTES_macOS.md)
- [INSTALL_NOTES_Windows.md](INSTALL_NOTES_Windows.md)

Installing from source
----------------------

### Required software


### Steps

1. git clone https://github.com/rsdoiel/opml
2. Change directory into the `opml` directory
3. Make to build, test and install

~~~shell
git clone https://github.com/rsdoiel/opml
cd opml
make
make test
make install
~~~

