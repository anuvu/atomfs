load helpers

function setup() {
    skopeo --insecure-policy copy docker://centos:latest oci:${TEST_DIR}/oci:centos
}

function teardown() {
    cleanup
    umount "${TEST_DIR}/centos" || true
    rm -rf "${TEST_DIR}/oci" "${TEST_DIR}/out.dump"
}

@test "import oci" {
    atomfs slurp-oci "${TEST_DIR}/oci"
    atomfs dump-db > "${TEST_DIR}/out.dump"
    grep centos "${TEST_DIR}/out.dump"
}

