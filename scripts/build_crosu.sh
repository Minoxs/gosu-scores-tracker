#!/usr/bin/env bash
ROOT=$(pwd)
SCRIPT_PATH=$(dirname $0)

cd "$SCRIPT_PATH" || exit
cd ../deps/crosu-pp || exit
cargo build --release

BUILT_LIB=./target/release/crosu_pp
LIB_DLL=$BUILT_LIB.dll
LIB_LINUX=$BUILT_LIB.so

[ -f $LIB_DLL ] && mv $LIB_DLL "$ROOT"
[ -f $LIB_LINUX ] && mv $LIB_LINUX "$ROOT"
