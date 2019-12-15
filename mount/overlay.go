package mount

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Mount struct {
	Source string
	Target string
	FSType string
	Opts   []string
}

type Mounts []Mount

func (ms Mounts) IsMountpoint(p string) bool {
	for _, m := range ms {
		if m.Target == p {
			return true
		}
	}

	return false
}

func ParseMounts() (Mounts, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't open /proc/self/mountinfo")
	}
	defer f.Close()

	mounts := []Mount{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		mount := Mount{}
		mount.Target = fields[4]

		for i := 5; i < len(fields); i++ {
			if fields[i] != "-" {
				continue
			}

			mount.FSType = fields[i+1]
			mount.Source = fields[i+2]
			mount.Opts = strings.Split(fields[i+3], ",")
		}

		mounts = append(mounts, mount)
	}

	return mounts, nil
}

func GetOverlayDirs(m Mount) []string {
	for _, opt := range m.Opts {
		if !strings.HasPrefix(opt, "lowerdir=") {
			continue
		}

		return strings.Split(strings.TrimPrefix(opt, "lowerdir="), ":")
	}

	return []string{}
}

