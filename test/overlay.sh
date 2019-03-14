#!/bin/bash -e

set -x

RESULT=failure

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
    set +x
    umount target || true
    rm -rf target work upper || true
    for i in $(seq ${levels}); do
        rm -rf layer${i}* || true
    done

    echo ==================
    echo RESULT: $RESULT
    echo ==================
}

trap cleanup EXIT HUP INT TERM

# do some custom setup to test whiteouts in lowerdirs
touch layer2-5091fe3751e44d94ca306eefc31888b877cf818a9f09a30d46cc09c32c03299a/hiddenfile
mknod layer1-dc14934a7f66ff36aaa4de4894679177ab787cb81d8f38ddb40f6de9cd6ba58b/hiddenfile c 0 0

mkdir layer1-dc14934a7f66ff36aaa4de4894679177ab787cb81d8f38ddb40f6de9cd6ba58b/hiddendir
mkdir layer2-5091fe3751e44d94ca306eefc31888b877cf818a9f09a30d46cc09c32c03299a/hiddendir
setfattr -n trusted.overlay.opaque -v y layer1-dc14934a7f66ff36aaa4de4894679177ab787cb81d8f38ddb40f6de9cd6ba58b/hiddendir

mknod layer1-dc14934a7f66ff36aaa4de4894679177ab787cb81d8f38ddb40f6de9cd6ba58b/hiddendir2 c 0 0
mkdir layer2-5091fe3751e44d94ca306eefc31888b877cf818a9f09a30d46cc09c32c03299a/hiddendir2

mkdir -p target
mkdir -p work
mkdir -p upper
mount -t overlay overlay "-olowerdir=${lowerdirs},upperdir=upper,workdir=work" target

for i in $(seq ${levels}); do
    stat target/${i} || (echo "${i} missing" && exit 1)
done

# check our special whiteouts
[ ! -f target/hiddenfile ]
[ ! -d target/hiddendir2 ]
[ -d target/hiddendir ] # we expect opaque dirs aren't hidden in lowerdirs

touch target/writable

ls -alh target
ls -alh work
ls -alh upper
RESULT=success
