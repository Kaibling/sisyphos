package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"sisyphos/models"

	"golang.org/x/crypto/ssh"
)

type SSHConnector struct{}

func NewSSHConnector() *SSHConnector {
	return &SSHConnector{}
}

func (s *SSHConnector) Execute(cfg models.SSHConfig, command string) (string, error) {
	// ssh config
	// hostKeyCallback, err := knownhosts.New(".sshwalk/known_hosts")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	connectionString := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	authMethods := []ssh.AuthMethod{}

	if cfg.Password != "" {
		fmt.Println("try password")
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	if cfg.PrivateKey != "" {
		fmt.Println("try ssh key")
		privKey := []byte(cfg.PrivateKey)
		signer, err := ssh.ParsePrivateKey(privKey)
		if err != nil {
			return "", err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	if len(authMethods) == 0 {
		return "", errors.New("not authentication methods available")
	}
	config := &ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // hostKeyCallback,
	}
	fmt.Printf("trying cmd: %s on '%s'\n", command, connectionString)
	// connect to ssh server
	conn, err := ssh.Dial("tcp", connectionString, config)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var buff bytes.Buffer
	session.Stdout = &buff
	if err := session.Run(command); err != nil {
		return "", err
	}
	return buff.String(), nil
}
