package ssh

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strings"

	"sisyphos/models"

	"golang.org/x/crypto/ssh"
)

type SSHConnector struct{}

func NewSSHConnector() *SSHConnector {
	return &SSHConnector{}
}

func (s *SSHConnector) Execute(cfg models.SSHConfig, command string) (string, error) {
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
	kh := KnownHosts{key: cfg.KnownKey}
	config := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: authMethods,
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(), // hostKeyCallback,
		HostKeyCallback: kh.ValidateHostKey,
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

func (s *SSHConnector) ReadHostKey(host string, port int) (string, error) {
	kh := KnownHosts{}
	sshConfig := &ssh.ClientConfig{HostKeyCallback: kh.ReadHostKey}
	_, _ = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), sshConfig) // TODO ????
	return kh.key, nil
}

type KnownHosts struct {
	key string
}

func (kh *KnownHosts) ValidateHostKey(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	current := fmt.Sprintf("%s %s %s", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	if kh.key != current {
		return fmt.Errorf("hostkey missmatch: unknown %s", current)
	}
	return nil
}

func (kh *KnownHosts) ReadHostKey(dialAddr string, addr net.Addr, key ssh.PublicKey) error {
	kh.key = fmt.Sprintf("%s %s %s", strings.Split(dialAddr, ":")[0], key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
	return nil
}
