ROOT_DIR=$(git rev-parse --show-toplevel)
TEST_DIR="${ROOT_DIR}/test"

function atomfs {
    root=
    run "${ROOT_DIR}/atomfs" --debug --base-dir "${TEST_DIR}/dir" "$@"
    echo "$output"
    [ "$status" -eq 0 ]
}

function cleanup {
    for d in $(ls "${TEST_DIR}/dir/mounts"); do
        umount -l "${TEST_DIR}/dir/mounts/$d"
    done
    rm -rf "${TEST_DIR}/dir" || true
}
