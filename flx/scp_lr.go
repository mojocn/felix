package flx

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dejavuzhou/felix/model"
	"github.com/pkg/sftp"
)

func ScpLR(h *model.Machine, localPath, remotePath string) error {
	c, err := NewSftpClient(h)
	if err != nil {
		return err
	}
	defer c.Close()
	return lrCopy(localPath, remotePath, c)
}
func lrCopy(local, remote string, c *sftp.Client) error {
	info, err := os.Lstat(local)
	if err != nil {
		return err
	}
	//if !info.Mode().IsRegular() {
	//	return fmt.Errorf("not support irregular file")
	//}
	if info.Mode()&os.ModeSymlink != 0 {
		return lrCopyL(local, remote, c)
	}
	if info.IsDir() {
		return lrCopyD(local, remote, c)
	}
	return lrCopyF(local, remote, info, c)
}
func lrCopyL(local, remote string, c *sftp.Client) error {
	realLocal, err := os.Readlink(local)
	if err != nil {
		return err
	}
	return lrCopy(realLocal, remote, c)
}
func lrCopyD(local, remote string, c *sftp.Client) error {
	contents, err := ioutil.ReadDir(local)
	if err != nil {
		return fmt.Errorf("ioutil read local dir failed %s", err)
	}
	for _, content := range contents {
		cs, cd := filepath.Join(local, content.Name()), filepath.Join(remote, content.Name())
		if err := lrCopy(cs, cd, c); err != nil {
			return fmt.Errorf("%s %s %s", err, cs, cd)
		}
	}
	return nil
}
func lrCopyF(local, remote string, info os.FileInfo, c *sftp.Client) error {
	localFile, err := os.Open(local)
	if err != nil {
		return fmt.Errorf("BrowserOpen local file failed %s", err)
	}
	defer localFile.Close()
	err = c.MkdirAll(toUnixPath(filepath.Dir(remote)))
	if err != nil {
		return fmt.Errorf("scp mkdir all failed %s", err)
	}
	remoteFile, err := c.Create(toUnixPath(remote))
	if err != nil {
		return fmt.Errorf("create remote file failed %s:%s", remote, err)
	}
	defer remoteFile.Close()
	err = c.Chmod(remoteFile.Name(), info.Mode())
	if err != nil {
		return fmt.Errorf("scp chmod failed %s", err)
	}
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("io copy failed %s", err)
	}
	return nil
}
