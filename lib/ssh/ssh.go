package ssh

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SSHConnector struct{}

func NewSSHConnector() *SSHConnector {
	return &SSHConnector{}
}

func (s *SSHConnector) Execute(connectionString, user, password, command string) (string, error) {
	// ssh config
	// hostKeyCallback, err := knownhosts.New(".sshwalk/known_hosts")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Printf("trying cmd: %s on '%s'\n", command, connectionString)
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // hostKeyCallback,
	}
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
