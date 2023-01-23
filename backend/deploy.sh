#!/bin/bash

CUR_PATH=`pwd`
echo $CUR_PATH
BUILD_DIR="$CUR_PATH/build"
BUILD_BIN_FILE="github-contributions"
PACKAGE_FILE="function.zip"

echo "Building ..."
if [[ ! -d $BUILD_DIR ]]; then
    mkdir $BUILD_DIR
fi
if [[ -e $BUILD_DIR/$BUILD_BIN_FILE ]]; then
    rm $BUILD_DIR/$BUILD_BIN_FILE
fi
GOOS=linux go build -o $BUILD_DIR/$BUILD_BIN_FILE
if [[ ! -e $BUILD_DIR/$BUILD_BIN_FILE ]]; then
    echo "Build FAILED!"
    exit -1
fi

echo "Packaging ..."
pushd $BUILD_DIR > /dev/null
if [[ -e $PACKAGE_FILE ]]; then
    rm $PACKAGE_FILE
fi
zip $PACKAGE_FILE $BUILD_BIN_FILE
echo "Deploying ..."
aws lambda update-function-code --function-name contributions --zip-file fileb://$PACKAGE_FILE > /dev/null
popd > /dev/null
