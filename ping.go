package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const DefaultMaxNanoSeconds = 1000 * 1000 * 1000

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

var (
	icmp    ICMP
	laddr   = net.IPAddr{IP: net.ParseIP("ip")}
	num     = 10
	timeout = 1000
	size    = 32
	stop    bool
)

func pingTtl(ip string) (int, int, int) {

	conn, err := net.DialTimeout("ip4:icmp", ip, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		fmt.Println(err)
		return 0, 0, DefaultMaxNanoSeconds
	}

	defer conn.Close()
	//icmp header
	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 1
	icmp.SequenceNum = 1

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	data := make([]byte, size)
	buffer.Write(data)
	data = buffer.Bytes()

	var SuccessTimes int
	var minTime int = DefaultMaxNanoSeconds
	var maxTime = -1
	var totalTime int
	for i := 0; i < num; i++ {
		icmp.SequenceNum = uint16(1)

		data[2] = byte(0)
		data[3] = byte(0)

		data[6] = byte(icmp.SequenceNum >> 8)
		data[7] = byte(icmp.SequenceNum)
		icmp.Checksum = CheckSum(data)
		data[2] = byte(icmp.Checksum >> 8)
		data[3] = byte(icmp.Checksum)

		t1 := time.Now()
		err := conn.SetDeadline(t1.Add(time.Duration(time.Duration(timeout) * time.Millisecond)))
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = conn.Write(data)
		if err != nil {
			log.Println(err)
			continue
		}
		buf := make([]byte, 65535)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		et := int(time.Since(t1).Nanoseconds())
		if minTime > et {
			minTime = et
		}
		if maxTime < et {
			maxTime = et
		}
		totalTime += et
		SuccessTimes++
		time.Sleep(10 * time.Millisecond)
	}
	if maxTime < 0 || SuccessTimes == 0 {
		return 0, 0, DefaultMaxNanoSeconds
	}

	return minTime, maxTime, totalTime / SuccessTimes
}

func CheckSum(data []byte) uint16 {
	var sum uint32
	var length = len(data)
	var index int
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	sum = uint32(sum>>16) + uint32(sum)
	sum = uint32(sum>>16) + uint32(sum)
	return uint16(^sum)
}
