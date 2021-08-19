package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func ReplaceHosts(old string, content string) string {
	var startText = "### start dnsfix"
	var endText = "### end dnsfix"
	var start = strings.Index(old, startText)
	var end = strings.Index(old, endText)
	if start <= 0 || end <= 0 {
		return fmt.Sprintf("%s\n%s\n%s\n%s\n", old, startText, content, endText)
	}
	var oldContent = old[start:end]
	new := strings.ReplaceAll(old, oldContent, startText+"\n"+content)
	fmt.Println(old)
	return new
}

func flushDns() error {
	var cmd *exec.Cmd
	switch goos {
	case "windows":
		cmd = exec.Command("ipconfig", "/flushdns")
	case "linux":
		cmd = exec.Command("service", "network", "restart")
	case "darwin":
		cmd = exec.Command("killall", "-HUP", "mDNSResponder")
	default:
		return errors.New("unknown os, ignore dns flush")
	}

	buf, err := cmd.CombinedOutput()

	fmt.Printf("%s \n", buf)
	return err
}
