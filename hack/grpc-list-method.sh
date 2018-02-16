#!/bin/sh

SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")
BASE_DIR=$SCRIPTPATH/..
SVC_PROTO=$BASE_DIR/pb/profile.proto

grep -oP "rpc \K(.*)\(.*\) returns" $SVC_PROTO | sed -n -e 's/^\([[:alnum:]]\+\).*/\1/p'
