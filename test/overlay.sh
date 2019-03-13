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

# do some custom setup to test whiteouts in lowerdirs
touch layer1-dc14934a7f66ff36aaa4de4894679177ab787cb81d8f38ddb40f6de9cd6ba58b/hiddenfile
mknod layer2-5091fe3751e44d94ca306eefc31888b877cf818a9f09a30d46cc09c32c03299a/hiddenfile c 0 0

mkdir layer1-dc14934a7f66ff36aaa4de4894679177ab787cb81d8f38ddb40f6de9cd6ba58b/hiddendir
mkdir layer2-5091fe3751e44d94ca306eefc31888b877cf818a9f09a30d46cc09c32c03299a/hiddendir
setfattr -n trusted.overlay.opaque -v y layer2-5091fe3751e44d94ca306eefc31888b877cf818a9f09a30d46cc09c32c03299a/hiddendir

mkdir -p target
mkdir -p work
mkdir -p upper
mount -t overlay overlay "-olowerdir=${lowerdirs},upperdir=upper,workdir=work" target
rmdir upper

for i in $(seq ${levels}); do
    stat target/${i} || (echo "${i} missing" && exit 1)
done

# check our special whiteouts
[ ! -f target/hiddenfile ]
[ ! -d target/hiddendir ]

touch target/writable

ls -alh target
ls -alh work
ls -alh upper
stat work/writable
