load helpers

function setup() {
    skopeo --insecure-policy copy docker://centos:latest oci:${TEST_DIR}/oci:centos
}

function teardown() {
    cleanup
    umount "${TEST_DIR}/centos" || true
    rm -rf "${TEST_DIR}/oci" "${TEST_DIR}/centos"
}

@test "import oci" {
    atomfs slurp-oci "${TEST_DIR}/oci"
    mkdir "${TEST_DIR}/centos"
    atomfs mount centos "${TEST_DIR}/centos"
    ls "${TEST_DIR}/dir"
    grep atomfs /proc/self/mountinfo
    atomfs umount "${TEST_DIR}/centos"
    atomfs fsck
}
