module github.com/anuvu/atomfs

require (
	github.com/anuvu/stacker v0.4.1-0.20190607155931-46d8ec0d501e
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/openSUSE/umoci v0.4.4
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/urfave/cli v1.20.0
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5 // indirect
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	golang.org/x/sys v0.0.0-20190405154228-4b34438f7a67
)

replace github.com/vbatts/go-mtree v0.4.4 => github.com/vbatts/go-mtree v0.4.5-0.20190122034725-8b6de6073c1a

replace github.com/openSUSE/umoci => github.com/tych0/umoci v0.1.1-0.20190402232331-556620754fb1
