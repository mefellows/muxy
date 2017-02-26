#!/bin/bash

# Environmental variables
CUR_DIR=$(pwd)
INPUT="${CUR_DIR}/pkg/"
OUTPUT="${CUR_DIR}/dist"

mkdir -p ${OUTPUT}
for TARGET in $(find ${INPUT} -mindepth 1 -maxdepth 1 -type d); do
    ARCHIVE_NAME=$(basename ${TARGET})
    pushd ${TARGET}
    zip -r ${OUTPUT}/${ARCHIVE_NAME}.zip ./*
    popd
done

# Generate shasum
pushd ${OUTPUT}
shasum -a 256 * > ./SHASUMS
popd

ls -l ${OUTPUT}
