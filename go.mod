module github.com/anuvu/atomfs

require (
	github.com/anuvu/stacker v0.5.1-0.20200212152300-22fec51333e7
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/openSUSE/umoci v0.4.6-0.20200206004913-cc1b6b2e346e
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/schollz/sqlite3dump v1.2.4
	github.com/urfave/cli v1.22.1
	golang.org/x/sys v0.0.0-20190913121621-c3b328c6e5a7
)

replace github.com/vbatts/go-mtree v0.4.4 => github.com/vbatts/go-mtree v0.4.5-0.20190122034725-8b6de6073c1a

go 1.13
