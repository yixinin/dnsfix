package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func ReplaceHosts(old string, content string) string {
	var startText = "### start speedup"
	var endText = "### end speedup"
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
	cmd := exec.Command("ipconfig", "/flushdns")
	buf, err := cmd.CombinedOutput()

	fmt.Printf("%s \n", buf)
	return err
}
