package internal

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type Target interface {
	Pre() (string, string, error)
	Post()
	String() string
}

var _ Target = (*DirectoryTarget)(nil)

const DirectoryTargetType = "directory"

type DirectoryTarget struct {
	From string `yaml:"from"`
}

func (t *DirectoryTarget) Pre() (string, string, error) {
	return t.From, t.From, nil
}

func (t *DirectoryTarget) Post() {
}

func (t *DirectoryTarget) String() string {
	return fmt.Sprintf("directory %s", t.From)
}

const zfsSnapshotMountBaseDir = "/restic-plus"

var _ Target = (*ZFSZvolTarget)(nil)

const ZFSDatasetTargetType = "zfs-dataset"

type ZFSDatasetTarget struct {
	From       string `yaml:"from"`
	snapshotId string
	mountDir   string
}

func (t *ZFSDatasetTarget) Pre() (string, string, error) {
	randomId := GenerateRandomString(12)
	t.snapshotId = t.From + "@" + randomId
	t.mountDir = path.Join(zfsSnapshotMountBaseDir, "zfs-dataset", t.From)

	if _, _, err := ExecCommandRetry("zfs", "snapshot", t.snapshotId); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(t.mountDir, 0o755); err != nil {
		return "", "", err
	}
	if _, _, err := ExecCommandRetry("mount", "-o", "ro", "-t", "zfs", t.snapshotId, t.mountDir); err != nil {
		return "", "", err
	}

	return t.From, t.mountDir, nil
}

func (t *ZFSDatasetTarget) Post() {
	if _, _, err := ExecCommandRetry("umount", t.mountDir); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSDatasetTargetType, err)
	}
	if err := os.RemoveAll(t.mountDir); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSDatasetTargetType, err)
	}
	if _, _, err := ExecCommandRetry("zfs", "destroy", t.snapshotId); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSDatasetTargetType, err)
	}
}

func (t *ZFSDatasetTarget) String() string {
	return fmt.Sprintf("zfs dataset %s", t.From)
}

var _ Target = (*ZFSZvolTarget)(nil)

const ZFSZvolTargetType = "zfs-zvol"

type ZFSZvolTarget struct {
	From        string `yaml:"from"`
	snapshotId  string
	zvolCloneId string
	zvolDevDir  string
	mountDir    string
}

func (t *ZFSZvolTarget) Pre() (string, string, error) {
	randomId := GenerateRandomString(12)
	t.snapshotId = t.From + "@" + randomId
	t.zvolCloneId = t.From + "-" + randomId
	t.zvolDevDir = path.Join("/dev/zvol", t.zvolCloneId)
	t.mountDir = path.Join(zfsSnapshotMountBaseDir, "zfs-zvol", t.From)

	if _, _, err := ExecCommandRetry("zfs", "snapshot", t.snapshotId); err != nil {
		return "", "", err
	}
	if _, _, err := ExecCommandRetry("zfs", "clone", t.snapshotId, t.zvolCloneId); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(t.mountDir, 0o755); err != nil {
		return "", "", err
	}
	if _, _, err := ExecCommandRetry("mount", "-o", "ro", t.zvolDevDir, t.mountDir); err != nil {
		return "", "", err
	}

	return t.From, t.mountDir, nil
}

func (t *ZFSZvolTarget) Post() {
	if _, _, err := ExecCommandRetry("umount", t.mountDir); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSZvolTargetType, err)
	}
	if err := os.RemoveAll(t.mountDir); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSZvolTargetType, err)
	}
	if _, _, err := ExecCommandRetry("zfs", "destroy", t.zvolCloneId); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSZvolTargetType, err)
	}
	if _, _, err := ExecCommandRetry("zfs", "destroy", t.snapshotId); err != nil {
		LogWarn.Printf("post for target %s failed: %v", ZFSZvolTargetType, err)
	}
}

func (t *ZFSZvolTarget) String() string {
	return fmt.Sprintf("zfs zvol %s", t.From)
}

func zfsGetMountPoint(id string) (string, error) {
	output, _, err := ExecCommandRetry("zfs", "get", "-o", "value", "mountpoint", id)
	if err != nil {
		return "", fmt.Errorf("unable to detect zfs mountpoint: %w", err)
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output from zfs command")
	}
	return strings.Trim(lines[1], " \t"), nil
}
