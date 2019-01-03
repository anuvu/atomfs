// Note: this file is unused right now, but it describes how we could
// auto-correct our mountns. The problem is with the `atomfs mount` command, we
// need to propagate the mount to the right place, back outside of the atomfs
// mountns. Need better ideas on how to do this.
package atomfs

/*
#define _GNU_SOURCE
#include <stdio.h>
#include <syscall.h>
#include <sched.h>
#include <linux/kcmp.h>

int create_atomfs_mntns(char *base_path)
{
	int ret;
	char buf[PATH_MAX];

	if (unshare(CLONE_NEWNS) < 0) {
		perror("unshare");
		return -1;
	}

	if (mkdir(base_path, 0755) < 0 && errno != EEXIST) {
		perror("mkdir");
		return -1;
	}

	snprintf(path, sizeof(path), "%s/ns", base_dir);
	ret = open(path, O_CREAT | O_WRONLY);
	if (ret < 0) {
		perror("create ns mount target");
		return -1;
	}
	close(ret);

	if (mount("/proc/self/ns/mnt", path, NULL, MS_BIND, NULL) < 0) {
		perror("mount");
		return -1;
	}

	return 0;
}

__attribute__((constructor)) void fixup_mntns(void)
{
	int ret, size, my_mnts, atomfs_mnts;
	char buf[4096], path[PATH_MAX], *base_dir = NULL;

	snprintf(buf, sizeof(buf), "/proc/self/ns/mnt")

	ret = open("/proc/self/cmdline", O_RDONLY);
	if (ret < 0) {
		perror("error: open");
		exit(96);
	}

	if ((size = read(ret, buf, sizeof(buf)-1)) < 0) {
		close(ret);
		perror("error: read");
		exit(96);
	}
	close(ret);

#define ADVANCE_ARG		\
	do {			\
		while (*cur) {	\
			cur++;	\
		}		\
		cur++;		\
	} while (0)

	// skip to base-dir
	while (1) {
		ADVANCE_ARG;

		if (cur - buf >= size) {
			base_dir = "/var/lib/atomfs";
			break;
		}

		if (strcmp(cur, "--base-dir"))
			continue;

		// skip --base-dir itself, and if the next arg is still valid,
		// use it as base-dir.
		ADVANCE_ARG;
		if (cur - buf >= size)
			base_dir = cur;
		else
			base_dir = "/var/lib/atomfs";
		break;
	}

	my_mntns = open("/proc/self/ns/mnt", O_RDONLY);
	if (my_mntns < 0) {
		perror("opening my mntns");
		exit(96);
	}

	snprintf(path, sizeof(path), "%s/ns", base_dir);
	atomfs_mntns = open(path, O_RDONLY);
	if (atomfs_mntns < 0) {
		if (errno == ENOENT) {
			if (create_atomfs_mntns(base_dir) < 0)
				exit(96)
			return;
		} else {
			perror("opening atomfs mntns");
			exit(96);
		}
	}

	ret = syscall(__NR_kcmp, getpid(), getpid(), KCMP_FILE, my_mntns, atomfs_mntns);
	close(my_mntns);
	close(atomfs_mntns);
	if (ret < 0) {
		perror("kcmp failed");
		exit(96);
	}

	// the namespaces match
	if (!ret)
		return;

	if (create_atomfs_mntns(base_dir) < 0)
		exit(96);
}
*/
