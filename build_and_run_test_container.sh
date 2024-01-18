#!/bin/bash

docker build -f $1 -t temp-dlv-image .

docker run temp-dlv-image
