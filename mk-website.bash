#!/bin/bash

PROJECT=opml

VERSION=$(grep -m 1 'Version =' opml.go | cut -d\" -f 2)

MOTTO="$PROJECT: A OPML parser package and opml sort utility"

function checkApp() {
    APP_NAME=$(which $1)
    if [ "$APP_NAME" = "" ] && [ ! -f "./bin/$1" ]; then
        echo "Missing $APP_NAME"
        exit 1
    fi
}

function softwareCheck() {
    for APP_NAME in $@; do
        checkApp $APP_NAME
    done
}

function MakePage () {
    motto="$1"
    version="$2"
    nav="$3"
    content="$4"
    html="$5"
    # Always use the latest compiled mkpage
    APP=$(which mkpage)
    if [ -f ./bin/mkpage ]; then
        APP="./bin/mkpage"
    fi

    echo "Rendering $html"
    $APP \
	    "title=text:$motto" \
        "version=text:$version" \
        "nav=$nav" \
        "content=$content" \
        page.tmpl > $html
}

echo "Checking necessary software is installed"
softwareCheck mkpage grep cut
echo "Generating website index.html"
MakePage "$MOTTO" "$PROJECT $VERSION" nav.md README.md index.html
echo "Generating install.html"
MakePage "$MOTTO" "$PROJECT $VERSION" nav.md INSTALL.md install.html
echo "Generating license.html"
MakePage "$MOTTO" "$PROJECT $VERSISON" nav.md "markdown:$(cat LICENSE)" license.html
