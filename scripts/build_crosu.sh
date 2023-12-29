#!/usr/bin/env bash
ROOT=$(pwd)
SCRIPT_PATH=$(dirname $0)

cd "$SCRIPT_PATH" || exit
cd ../deps/crosu-pp || exit
cargo build --release

LIB_PATH=./target/release
LIB_NAME=crosu_pp
LIB_DLL=$LIB_PATH/$LIB_NAME.dll
LIB_LINUX=$LIB_PATH/lib$LIB_NAME.so

[ -f $LIB_DLL ] && mv $LIB_DLL "$ROOT"
[ -f $LIB_LINUX ] && mv $LIB_LINUX "$ROOT"
