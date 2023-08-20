package util

import (
	"DouYin/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

var SftpClient *sftp.Client
var SshClient *ssh.Client

func InitSFTP() {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(config.FtpPsw))
	clientConfig = &ssh.ClientConfig{
		User:            config.FtpUser,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	if config.TypeSSH == "password" {
		clientConfig.Auth = []ssh.AuthMethod{ssh.Password(config.PasswordSSH)} // Use password-based authentication if configured
	}
	// connet to ssh
	if SshClient, err = ssh.Dial("tcp", config.FtpAdd, clientConfig); err != nil {
		log.Fatal(" bulid ssh client error", err)
	}
	// create sftp client
	if SftpClient, err = sftp.NewClient(SshClient); err != nil {
		log.Fatal(" bulid ssh client error", err)
	}
	keepsshAlive()
}
func keepsshAlive() {
	// Keep SSH session alive for FFmpeg processing
	time.Sleep(time.Duration(config.SSHHeartbeatTime) * time.Second)
	session, _ := SshClient.NewSession()
	session.Close()
}
