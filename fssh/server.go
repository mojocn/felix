package fssh

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/mojocn/felix/util"
	"github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var hostKeySigner gossh.Signer

func LoadOrCreateKey() {
	s, err := createOrLoadKeySigner()
	if err != nil {
		log.Fatal(err)
	}
	hostKeySigner = s
}

func Run() {
	s := &ssh.Server{
		Addr:            ":88",
		Handler:         helloHandler,
		PasswordHandler: passwordH,
	}
	s.AddHostKey(hostKeySigner)
	log.Fatal(s.ListenAndServe())

}

func passwordH(ctx ssh.Context, password string) bool {
	user := ctx.User()
	return user == "felix"
}

func helloHandler(s ssh.Session) {
	io.WriteString(s, fmt.Sprintf("你好欢迎来到felix堡垒机 自定义SSH 服务 %s\n", s.User()))

	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		io.WriteString(s, "不是PTY请求.\n")
		s.Exit(1)
	}
	sshConf, err := util.NewSshClientConfig("pi", "", "password", "", "")
	if err != nil {
		io.WriteString(s, err.Error())
		return
	}
	// Connect to ssh server
	conn, err := gossh.Dial("tcp", "home.mojotv.cn:22", sshConf)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()
	// CreateUserOfRole a fss
	fss, err := conn.NewSession()
	if err != nil {
		log.Fatal("unable to create fss: ", err)
	}
	defer fss.Close()
	// Set up terminal modes
	modes := gossh.TerminalModes{
		gossh.ECHO:          1,     // disable echoing
		gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := fss.RequestPty(ptyReq.Term, ptyReq.Window.Height, ptyReq.Window.Width, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	go func() {
		for win := range winCh {
			err := fss.WindowChange(win.Height, win.Width)
			if err != nil {
				logrus.WithError(err).Error("windows size changed")
			}
		}
	}()
	//stdinP, err := fss.StdinPipe()
	//if err != nil {
	//	logrus.WithError(err).Error("stdin Pipe")
	//	return
	//}
	//stdoutP, err := fss.StdoutPipe()
	//if err != nil {
	//	logrus.WithError(err).Error("stdot Pipe")
	//	return
	//}
	//stderrP, err := fss.StderrPipe()
	//if err != nil {
	//	logrus.WithError(err).Error("stderr Pipe")
	//	return
	//}
	//
	//go io.Copy(s,stderrP)
	//go io.Copy(s,stdoutP)
	//go io.Copy(stdinP, s) // stdin
	// Start remote shell
	fss.Stderr = s
	fss.Stdin = s
	fss.Stdout = s
	if err := fss.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}
	fss.Wait()
}

func createOrLoadKeySigner() (gossh.Signer, error) {
	keyPath := filepath.Join(".", "fssh.rsa")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		stderr, err := exec.Command("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "").CombinedOutput()
		output := string(stderr)
		if err != nil {
			return nil, fmt.Errorf("Fail to generate private key: %v - %s", err, output)
		}
	}
	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return gossh.ParsePrivateKey(privateBytes)
}
