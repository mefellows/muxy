#! /bin/bash -xe

cd /opt/graphite/storage
for f in $(find $1 -iname "*.wsp")
do
    if [ -a $f ]
    then
       echo "resizing $f" 
       /opt/graphite/bin/whisper-resize.py $f 10s:6h,1min:6d,10min:1800d
    fi
done

