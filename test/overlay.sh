#!/bin/bash -e

set -x

levels=10

lowerdirs=""
for i in $(seq ${levels}); do
    name=layer${i}
    # simulate the layer name being a sha256 sum to make sure mntopts has
    # enough room.
    name="${name}-$(echo $name | sha256sum | cut -f1 -d" ")"
    mkdir -p $name
    touch $name/${i}
    lowerdirs="${lowerdirs}${name}:"
done
lowerdirs=${lowerdirs::-1}

function cleanup() {
    umount target || true
    rm -rf target work upper || true
    for i in $(seq ${levels}); do
        rm -rf layer${i}* || true
    done
}

trap cleanup EXIT HUP INT TERM

mkdir -p target
mkdir -p work
mkdir -p upper
mount -t overlay overlay "-olowerdir=${lowerdirs},upperdir=upper,workdir=work" target
rmdir upper

for i in $(seq ${levels}); do
    stat target/${i} >& /dev/null || (echo "${i} missing" && exit 1)
done

touch target/writable

ls -alh target
ls -alh work
ls -alh upper
stat work/writable
