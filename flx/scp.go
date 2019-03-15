package flx

import (
	"path/filepath"

	"github.com/dejavuzhou/felix/model"
	"github.com/pkg/sftp"
)

const maxPacket = 1 << 15

func NewSftpClient(h *model.Machine) (*sftp.Client, error) {
	conn, err := NewSshClient(h)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(conn, sftp.MaxPacket(maxPacket))
}
func toUnixPath(path string) string {
	return filepath.Clean(path)
}
