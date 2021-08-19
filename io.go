package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Config struct {
	Dnss    [][]string `json:"dnss"`
	Domains []string   `json:"domains"`
}

func readHosts() string {
	var path string
	switch goos {
	case "windows":
		path = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux", "darwin":
		path = "/etc/hosts"
	default:
		return ""
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Printf("read %s buffer error:%v\n", path, err)
		return ""
	}
	return string(buf)
}

func saveHosts(content string) error {
	var path string
	switch goos {
	case "windows":
		path = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux", "darwin":
		path = "/etc/hosts"
	default:
		return errors.New("unknown os")
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		log.Printf("open %s buffer error:%v", path, err)
		return err
	}
	_, err = io.WriteString(f, content)
	if err != nil {
		log.Printf("write error:%v\n", err)
	}
	return err
}

func readConfig(path string) Config {
	var config Config
	buf, err := os.ReadFile(path)
	if err != nil {
		panic("no config file: " + path)
	}

	err = json.Unmarshal(buf, &config)
	if err != nil {
		log.Println("read dns.txt buffer error:", err)
		return config
	}
	return config
}

func readPid() string {
	buf, err := os.ReadFile("pid.txt")
	if err != nil {
		log.Println("read dns.txt buffer error:", err)
		return ""
	}
	return string(buf)
}

func savePid(pid string) {
	log.Println(os.WriteFile("pid.txt", []byte(pid), os.ModePerm))
}
