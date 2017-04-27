package main

import (
	"log"
	"os"
	"time"

	"github.com/ziutek/telnet"
)

const timeout = 10 * time.Second

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func expect(t *telnet.Conn, d ...string) {
	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}

func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}

func main() {
	if len(os.Args) != 3 {
		log.Printf("Usage: %s command param ", os.Args[0])
		return
	}
	command := os.Args[1]
	params := os.Args[2]

	t, err := telnet.Dial("tcp", "192.168.1.206:10002")
	checkErr(err)
	t.SetUnixWriteMode(true)
	var data []byte
	switch command {
	case "power":
		if params == "off" {
			sendln(t, "POWR0   ")
		} else {
			sendln(t, "POWR1   ")
		}
		data, err = t.ReadBytes('\r')

	case "volume":
		vol := rightPad(params)
		sendln(t, "VOLM"+vol)
		data, err = t.ReadBytes('\r')
	default:
		log.Fatalln("bad command: " + command)
	}
	checkErr(err)
	os.Stdout.Write(data)
	os.Stdout.WriteString("\n")
}

func rightPad(s string) string {
	switch len(s) {
	case 1:
		return s + "   "
	case 2:
		return s + "  "
	case 3:
		return s + " "
	default:
		return "10  "
	}
}
