package internal

import (
	"compress/bzip2"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
)

var (
	resticDir        = path.Join(os.TempDir(), "restic")
	sftpIdentityFile = path.Join(resticDir, "sftp-identity.key")
)

func ExecRestic(ctx *Context, args ...string) error {
	config := ctx.Config

	resticRepository := fmt.Sprintf("sftp://%s@%s:%d/", config.SFTP.User, config.SFTP.Host, config.SFTP.Port)
	resticPassword := config.Restic.Password

	sshCommand := fmt.Sprintf("ssh -i %s -p %d %s@%s -s sftp", sftpIdentityFile, config.SFTP.Port, config.SFTP.User, config.SFTP.Host)
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs,
		"-o",
		fmt.Sprintf("sftp.command=%s", sshCommand),
	)
	cmdArgs = append(cmdArgs, args...)
	cmdEnv := []string{
		"RESTIC_REPOSITORY=" + resticRepository,
		"RESTIC_PASSWORD=" + resticPassword,
	}

	LogDebug.Printf("Executing restic command %v\n", cmdArgs)
	if _, _, err := ExecCommandWithOpts(ExecCommandOpts{
		Name:   ctx.ResticBinary,
		Args:   cmdArgs,
		Env:    cmdEnv,
		Logger: LogRestic,
	}); err != nil {
		return fmt.Errorf("restic command failed: %w", err)
	}
	return nil
}

func prepareResticBinary(version string, target string) error {
	resticBinUrl := fmt.Sprintf("https://github.com/restic/restic/releases/download/v%s/restic_%s_%s_%s.bz2", version, version, runtime.GOOS, runtime.GOARCH)

	if stat, err := os.Stat(target); err == nil && !stat.IsDir() {
		return nil
	}

	err := os.MkdirAll(path.Dir(target), 0o755)
	if err != nil {
		return err
	}
	LogDebug.Printf("Downloading %s to %s\n", resticBinUrl, target)
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(resticBinUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	archive := bzip2.NewReader(resp.Body)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, archive)
	if err != nil {
		return err
	}

	return nil
}
