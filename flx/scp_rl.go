package flx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dejavuzhou/felix/model"
	"github.com/pkg/sftp"
)

func ScpRL(h *model.Machine, remotePath, localPath string) error {
	c, err := NewSftpClient(h)
	if err != nil {
		return err
	}
	defer c.Close()
	return rlCopy(remotePath, localPath, c)
}

func rlCopy(remote, local string, c *sftp.Client) error {
	info, err := c.Lstat(toUnixPath(remote))
	if err != nil {
		return err
	}
	//if !info.Mode().IsRegular() {
	//	return fmt.Errorf("not support irregular file")
	//}
	if info.Mode()&os.ModeSymlink != 0 {
		return rlCopyL(local, remote, c)
	}
	if info.IsDir() {
		return rlCopyD(remote, local, c)
	}
	return rlCopyF(remote, local, info, c)
}
func rlCopyL(remote, local string, c *sftp.Client) error {
	realRemote, err := c.ReadLink(toUnixPath(remote))
	if err != nil {
		return err
	}
	return lrCopy(realRemote, local, c)
}
func rlCopyD(remote, local string, c *sftp.Client) error {
	contents, err := c.ReadDir(toUnixPath(remote))
	if err != nil {
		return fmt.Errorf("ioutil read scp remote dir failed %s", err)
	}
	for _, info := range contents {
		cdL, csR := filepath.Join(local, info.Name()), filepath.Join(remote, info.Name())
		//mkdir local dir by remote
		err := os.MkdirAll(filepath.Dir(cdL), info.Mode())
		if err != nil {
			return fmt.Errorf("os local sub mkdir all failed,%s", err)
		}
		csR = toUnixPath(csR)
		if err := rlCopy(csR, cdL, c); err != nil {
			return fmt.Errorf("dir walk remote:%s, local:%s, %s", csR, cdL, err)
		}
	}
	return nil
}
func rlCopyF(remote, local string, info os.FileInfo, c *sftp.Client) error {
	rFile, err := c.Open(toUnixPath(remote))
	if err != nil {
		return fmt.Errorf("BrowserOpen scp remote file failed %s", err)
	}
	defer rFile.Close()

	lFile, err := os.Create(local)
	if err != nil {
		return fmt.Errorf("os create local file failed:%s %s", local, err)
	}
	defer lFile.Close()

	size, err := io.Copy(lFile, rFile)
	if err != nil {
		return fmt.Errorf("io copy remote to local failed.size:%d %s", size, err)
	}

	err = os.Chmod(lFile.Name(), info.Mode())
	if err != nil {
		return fmt.Errorf("os local chmod failed %s", err)
	}
	return nil
}
