#!/bin/bash

PROJECT="$1"

MOTTO="$2"

VERSION="$3"

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

if [ "$PROJECT" = "" ] || [ "$MOTTO" = "" ] || [ "$VERSION" = "" ]; then
    echo "USAGE: mk-website.absh PROJECT MOTTO VERSION"
    exit 1
fi
echo "Checking necessary software is installed"
softwareCheck mkpage grep cut
echo "Generating website index.html"
MakePage "$MOTTO" "$PROJECT $VERSION" nav.md README.md index.html
echo "Generating install.html"
MakePage "$MOTTO" "$PROJECT $VERSION" nav.md INSTALL.md install.html
echo "Generating license.html"
MakePage "$MOTTO" "$PROJECT $VERSISON" nav.md "markdown:$(cat LICENSE)" license.html
