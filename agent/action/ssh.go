package action

import (
	"bufio"
	"fmt"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Hostname string
	Port     int
	Username string
	Password string
}

func NewSSHClient(hostname string, port int, username, password string) SSHClient {
	return SSHClient{
		Hostname: hostname,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *SSHClient) Run(cmd string) ([]string, error) {
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		//TODO: fix this
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Hostname, s.Port), config)

	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()

	outPipe, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.Run(cmd)
	var lines []string
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		reader := bufio.NewScanner(outPipe)
		for reader.Scan() {
			lines = append(lines, reader.Text())
		}
		wg.Done()
	}()

	wg.Wait()

	return lines, nil
}
