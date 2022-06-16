package internal

import (
	"fmt"
	"os"
	"path"
)

type Stanza interface {
	Pre() ([]string, error)
	Post()
	String() string
}

var _ Stanza = (*DirectoryStanza)(nil)

const DirectoryStanzaType = "directory"

type DirectoryStanza struct {
	From string `yaml:"from"`
}

func (t *DirectoryStanza) Pre() ([]string, error) {
	return []string{t.From}, nil
}

func (t *DirectoryStanza) Post() {
}

func (t *DirectoryStanza) String() string {
	return fmt.Sprintf("directory %s", t.From)
}

const zfsSnapshotMountBaseDir = "/restic-plus"

var _ Stanza = (*ZFSZvolStanza)(nil)

const ZFSDatasetStanzaType = "zfs-dataset"

type ZFSDatasetStanza struct {
	From       string `yaml:"from"`
	snapshotId string
	mountDir   string
}

func (t *ZFSDatasetStanza) Pre() ([]string, error) {
	randomId := GenerateRandomString(12)
	t.snapshotId = t.From + "@" + randomId
	t.mountDir = path.Join(zfsSnapshotMountBaseDir, "zfs-dataset", t.From)

	if _, _, err := ExecCommandRetry("zfs", "snapshot", t.snapshotId); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(t.mountDir, 0o755); err != nil {
		return nil, err
	}
	if _, _, err := ExecCommandRetry("mount", "-o", "ro", "-t", "zfs", t.snapshotId, t.mountDir); err != nil {
		return nil, err
	}

	return []string{t.mountDir}, nil
}

func (t *ZFSDatasetStanza) Post() {
	if _, _, err := ExecCommandRetry("umount", t.mountDir); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSDatasetStanzaType, err)
	}
	if err := os.RemoveAll(t.mountDir); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSDatasetStanzaType, err)
	}
	if _, _, err := ExecCommandRetry("zfs", "destroy", t.snapshotId); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSDatasetStanzaType, err)
	}
}

func (t *ZFSDatasetStanza) String() string {
	return fmt.Sprintf("zfs dataset %s", t.From)
}

var _ Stanza = (*ZFSZvolStanza)(nil)

const ZFSZvolStanzaType = "zfs-zvol"

type ZFSZvolStanza struct {
	From        string `yaml:"from"`
	snapshotId  string
	zvolCloneId string
	zvolDevDir  string
	mountDir    string
}

func (t *ZFSZvolStanza) Pre() ([]string, error) {
	randomId := GenerateRandomString(12)
	t.snapshotId = t.From + "@" + randomId
	t.zvolCloneId = t.From + "-" + randomId
	t.zvolDevDir = path.Join("/dev/zvol", t.zvolCloneId)
	t.mountDir = path.Join(zfsSnapshotMountBaseDir, "zfs-zvol", t.From)

	if _, _, err := ExecCommandRetry("zfs", "snapshot", t.snapshotId); err != nil {
		return nil, err
	}
	if _, _, err := ExecCommandRetry("zfs", "clone", t.snapshotId, t.zvolCloneId); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(t.mountDir, 0o755); err != nil {
		return nil, err
	}
	if _, _, err := ExecCommandRetry("mount", "-o", "ro", t.zvolDevDir, t.mountDir); err != nil {
		return nil, err
	}

	return []string{t.mountDir}, nil
}

func (t *ZFSZvolStanza) Post() {
	if _, _, err := ExecCommandRetry("umount", t.mountDir); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSZvolStanzaType, err)
	}
	if err := os.RemoveAll(t.mountDir); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSZvolStanzaType, err)
	}
	if _, _, err := ExecCommandRetry("zfs", "destroy", t.zvolCloneId); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSZvolStanzaType, err)
	}
	if _, _, err := ExecCommandRetry("zfs", "destroy", t.snapshotId); err != nil {
		LogWarn.Printf("post for stanza %s failed: %v", ZFSZvolStanzaType, err)
	}
}

func (t *ZFSZvolStanza) String() string {
	return fmt.Sprintf("zfs zvol %s", t.From)
}
