package utils

import "fmt"

// GetDockerLoginCmd returns docker login cli command
func GetDockerLoginCmd(username, password string) string {
	if username == "" || password == "" {
		return ""
	}
	return fmt.Sprintf("docker login -u %s -p %s", username, password)
}

// GetPodmanLoginCmd returns podman login cli command
func GetPodmanLoginCmd(username, password string) string {
	if username == "" || password == "" {
		return ""
	}
	return fmt.Sprintf("podman login -u %s -p %s", username, password)
}
