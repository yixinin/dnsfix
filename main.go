package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"

	"github.com/miekg/dns"
)

var configPath string

func main() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		configPath = flag.Args()[0]
	} else {
		configPath = "./config.json"
	}

	var config = readConfig(configPath)
	var ch = make(chan A)
	var wg sync.WaitGroup

	var as = make(map[string]ASlice, len(config.Domains))

	for _, ds := range config.Dnss {
		var dnsCopy = ds[0]
		for _, domain := range config.Domains {
			var domainCopy = domain
			wg.Add(1)
			go dnsQuery(ch, &wg, domainCopy, dnsCopy)
		}
	}

	go func() {
		for a := range ch {
			as[a.Domain] = append(as[a.Domain], a)
		}
	}()

	wg.Wait()
	close(ch)

	var localDns strings.Builder

	for addr, a := range as {
		if len(a) == 0 {
			continue
		}
		sort.Sort(a)
		localDns.WriteString(fmt.Sprintf("%s\t%s\t# %d\n", a[0].Ip, addr, a[0].Ttl))
	}
	hostsOld := readHosts(config.HostPaths)
	if hostsOld == "" {
		fmt.Println("hosts empty")
		return
	}
	hostsNew := ReplaceHosts(hostsOld, localDns.String())
	saveHosts(config.HostPaths, hostsNew)
	flushDns()
	fmt.Println("speedup success. press any key to exit...")

	go func() {
		_, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	time.Sleep(30 * time.Second)
}

type A struct {
	Domain string
	Ip     string
	Ttl    int
}

type ASlice []A

func (a ASlice) Len() int           { return len(a) }
func (a ASlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ASlice) Less(i, j int) bool { return a[i].Ttl < a[j].Ttl }

func dnsQuery(ch chan A, wg *sync.WaitGroup, domain string, dnsIp string) {
	defer func() {
		recover()
		wg.Done()
	}()
	var client = &dns.Client{
		Timeout: 5000 * time.Millisecond,
	}
	var msg = &dns.Msg{}
	msg.SetQuestion(domain+".", dns.TypeA)
	r, _, err := client.Exchange(msg, dnsIp+":53")
	if err != nil {
		return
	}
	for _, ans := range r.Answer {
		record, isType := ans.(*dns.A)
		if isType {
			ip := record.A.String()
			min, max, avg := pingTtl(ip)
			if avg == 0 {
				fmt.Println(min, max)
			}
			if avg == 1000*1000*1000 {
				return
			}

			ch <- A{
				Ip:     ip,
				Ttl:    avg,
				Domain: domain,
			}
		}
	}
}
