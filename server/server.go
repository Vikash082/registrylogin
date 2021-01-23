package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Vikash082/registrylogin/utils"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

const (
	// ClientTimeout in seconds
	ClientTimeout = "60s"
)

// Server represent server
type Server struct {
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	IP             string `yaml:"ip"`
	Authfileloc    string `yaml:"authfileloc"`
	RegistryClient string `yaml:"registryclient"`
	Sudo           bool   `yaml:"sudo"`
}

// Execute given commands on remote server
func (s Server) Execute(command string) error {
	client, err := s.GetConnection()
	if err != nil {
		log.Fatalf("Issue getting client: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Issue getting session: %v", err)
	}
	defer session.Close()
	outputByte, err := session.CombinedOutput(command)
	if err != nil {
		return err
	}
	fmt.Println(string(outputByte))
	return nil
}

// GetConnection returns ssh client for the given server
func (s Server) GetConnection() (*ssh.Client, error) {
	tout, err := time.ParseDuration(ClientTimeout)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         tout,
	}
	if s.Password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(s.Password),
		}
	} else if s.Authfileloc != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(s.getSigner()),
		}
	} else {
		log.Fatalf("None of the authentication password or key provided for server: %s", s.IP)
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.IP, "22"), config)
}

func (s Server) getSigner() ssh.Signer {
	keyPath, err := homedir.Expand(s.Authfileloc)
	if err != nil {
		log.Fatalf("error getting the path %s on server %s : %v", s.Authfileloc, s.IP, err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	return signer
}

// GetLoginCommand returns the login command for registry cli
func (s Server) GetLoginCommand(dockerUser, dockerPassword, podmanUser, podmanPassword string) string {
	var cli string = s.RegistryClient
	switch cli {
	case "docker":
		return s.sudoRun(utils.GetDockerLoginCmd(dockerUser, dockerPassword))
	case "podman":
		return s.sudoRun(utils.GetPodmanLoginCmd(podmanUser, podmanPassword))
	default:
		return s.sudoRun(utils.GetDockerLoginCmd(dockerUser, dockerPassword))
	}
}

func (s Server) sudoRun(command string) string {
	if s.Sudo {
		return fmt.Sprintf("%s %s", "sudo", command)
	}
	return command
}
