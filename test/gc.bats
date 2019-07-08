load helpers

function teardown() {
    cleanup
}

@test "gc without adding atoms is ok" {
    atomfs init
    atomfs gc
}
