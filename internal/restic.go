package internal

import (
	"compress/bzip2"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
)

var (
	resticVersion = "0.13.1"
	resticBinUrl  = fmt.Sprintf("https://github.com/restic/restic/releases/download/v%s/restic_%s_%s_%s.bz2", resticVersion, resticVersion, runtime.GOOS, runtime.GOARCH)

	resticDir        = path.Join(os.TempDir(), "restic")
	resticBin        = path.Join(resticDir, "restic")
	sftpIdentityFile = path.Join(resticDir, "sftp-identity.key")
)

func prepareBinary() error {
	if stat, err := os.Stat(resticBin); err == nil && !stat.IsDir() {
		return nil
	}

	err := os.MkdirAll(resticDir, 0o755)
	if err != nil {
		return err
	}
	Debug.Printf("Downloading %s to %s\n", resticBinUrl, resticBin)
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
	file, err := os.OpenFile(resticBin, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
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

func prepareIdentity(identityPrivateKey string) error {
	err := os.MkdirAll(resticDir, 0o755)
	if err != nil {
		return err
	}
	ioutil.WriteFile(sftpIdentityFile, []byte(identityPrivateKey), 0o600)

	return nil
}

func Restic(args ...string) error {
	config := Config{}
	if err := config.LoadFromFile(""); err != nil {
		return err
	}

	if err := prepareBinary(); err != nil {
		return fmt.Errorf("unable to prepare underlying restic binary: %w", err)
	}

	if err := prepareIdentity(config.SFTP.IdentityPrivateKey); err != nil {
		return fmt.Errorf("unable to prepare identity file: %w", err)
	}

	resticRepository := fmt.Sprintf("sftp://%s@%s:%d/", config.SFTP.User, config.SFTP.Host, config.SFTP.Port)
	resticPassword := config.Restic.Password
	cmdEnv := os.Environ()
	cmdEnv = append(cmdEnv,
		"RESTIC_REPOSITORY="+resticRepository,
		"RESTIC_PASSWORD="+resticPassword,
	)

	sshCommand := fmt.Sprintf("ssh -i %s -p %d %s@%s -s sftp", sftpIdentityFile, config.SFTP.Port, config.SFTP.User, config.SFTP.Host)
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs,
		"-o",
		fmt.Sprintf("sftp.command=%s", sshCommand),
	)
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(resticBin, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = cmdEnv

	Debug.Printf("Executing restic command %v\n", cmdArgs)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restic command failed: %w", err)
	}
	return nil
}
