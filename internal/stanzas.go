package internal

import (
	"fmt"
	"os"
	"path"
	"strings"
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
	From  string `yaml:"from"`
	inner []struct {
		snapshotId  string
		zvolCloneId string
		zvolDevDir  string
		mountDir    string
	}
}

func (t *ZFSZvolStanza) Pre() ([]string, error) {
	zvols, err := zfsListZvols(t.From)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, zvol := range zvols {
		randomId := GenerateRandomString(12)

		t2 := struct {
			snapshotId  string
			zvolCloneId string
			zvolDevDir  string
			mountDir    string
		}{}

		t2.snapshotId = zvol + "@" + randomId
		t2.zvolCloneId = zvol + "-" + randomId
		t2.zvolDevDir = path.Join("/dev/zvol", t2.zvolCloneId)
		t2.mountDir = path.Join(zfsSnapshotMountBaseDir, "zfs-zvol", zvol)

		if _, _, err := ExecCommandRetry("zfs", "snapshot", t2.snapshotId); err != nil {
			return nil, err
		}
		if _, _, err := ExecCommandRetry("zfs", "clone", t2.snapshotId, t2.zvolCloneId); err != nil {
			return nil, err
		}
		if err := os.MkdirAll(t2.mountDir, 0o755); err != nil {
			return nil, err
		}
		if _, _, err := ExecCommandRetry("mount", "-o", "ro", t2.zvolDevDir, t2.mountDir); err != nil {
			return nil, err
		}
		result = append(result, t2.mountDir)
		t.inner = append(t.inner, t2)
	}

	return result, nil
}

func (t *ZFSZvolStanza) Post() {
	for _, t := range t.inner {
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
}

func (t *ZFSZvolStanza) String() string {
	return fmt.Sprintf("zfs zvols %s", t.From)
}

func zfsListZvols(parent string) ([]string, error) {
	output, _, err := ExecCommand("zfs", "list", "-t", "volume", "-o", "name", "-d", "256", parent)
	if err != nil {
		return nil, fmt.Errorf("unable to list zfs zvols: %w", err)
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected output from zfs command")
	}

	zvols := []string{}
	for _, line := range lines[1:] {
		zvol := strings.Trim(line, " ")
		if zvol != "" {
			zvols = append(zvols, zvol)
		}
	}
	return zvols, nil
}
