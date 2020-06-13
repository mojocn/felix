package fssh

import (
	"encoding/json"
	"errors"
	"github.com/libragen/felix/util"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"time"
)

type fsshToken struct {
	Uid uint      `json:"uid"`
	Mid uint      `json:"mid"`
	Ex  time.Time `json:"ex"`
}

func TokenGenerate(userId, machineId uint, ex time.Duration) (secret string, err error) {
	t := fsshToken{
		Uid: userId,
		Mid: machineId,
		Ex:  time.Now().Add(ex),
	}
	bs, err := json.Marshal(t)
	if err != nil {
		return
	}
	key := viper.GetString("app.secret")
	return util.AesEncrypt(bs, key)
}

func TokenToSession(token string) (c *ssh.Client, err error) {
	key := viper.GetString("app.secret")
	bs, err := util.AesDecrypt(token, key)
	if err != nil {
		return
	}
	t := fsshToken{}
	err = json.Unmarshal(bs, &t)
	if err != nil {
		return
	}
	if t.Ex.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	sshConf, err := util.NewSshClientConfig("pi", "000", "000", "", "")
	if err != nil {
		return nil, err
	}
	// Connect to ssh server
	return ssh.Dial("tcp", "home.mojotv.cn:22", sshConf)
}
