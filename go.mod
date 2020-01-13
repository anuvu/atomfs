module github.com/anuvu/atomfs

require (
	github.com/anuvu/stacker v0.4.1-0.20200113225657-ae7c48e3dc7f
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/openSUSE/umoci v0.4.4
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/schollz/sqlite3dump v1.2.4
	github.com/urfave/cli v1.20.0
	golang.org/x/sys v0.0.0-20190902133755-9109b7679e13
)

replace github.com/vbatts/go-mtree v0.4.4 => github.com/vbatts/go-mtree v0.4.5-0.20190122034725-8b6de6073c1a

replace github.com/openSUSE/umoci => github.com/tych0/umoci v0.1.1-0.20190402232331-556620754fb1

go 1.13
