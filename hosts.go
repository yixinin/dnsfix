package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func ReplaceHosts(old string, content string) string {
	var startText = "### start github"
	var endText = "### end github"
	var start = strings.Index(old, startText)
	var end = strings.Index(old, endText)
	if start <= 0 || end <= 0 {
		return old + "\n" + content
	}
	var oldContent = old[start:end]
	new := strings.ReplaceAll(old, oldContent, startText+"\n"+content)
	fmt.Println(old)
	return new
}

func flushDns() {
	cmd := exec.Command("ipconfig", "/flushdns")
	buf, err := cmd.CombinedOutput()
	fmt.Printf("%s %v", buf, err)
}
