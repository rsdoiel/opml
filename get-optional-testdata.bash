#!/bin/bash

START=$(pwd)
cd testdata
curl -L -O "http://scripting.com/misc/userlandSamples.zip"
cd "$START"
