package bchainlibs

import (
	"net"
	"strconv"
	"fmt"
	"github.com/op/go-logging"
	"encoding/json"
)

func SendToNetwork(serverIp string, serverPort string, channel <-chan string, toLog bool, log *logging.Logger, me net.IP) {
	Server, err := net.ResolveUDPAddr(Protocol, serverIp+serverPort)
	CheckError(err, log)
	Local, err := net.ResolveUDPAddr(Protocol, me.String()+LocalPort)
	CheckError(err, log)
	Conn, err := net.DialUDP(Protocol, Local, Server)
	CheckError(err, log)
	defer Conn.Close()

	for {
		j, more := <-channel
		if more {
			if Conn != nil {
				buf := []byte(j)
				_, err = Conn.Write(buf)
				if toLog {
					log.Debug(me.String() + " " + j + " MESSAGE_SIZE=" + strconv.Itoa(len(buf)))
					log.Debug(me.String() + " SENDING_MESSAGE=1")
				}
				CheckError(err, log)
			}
		} else {
			fmt.Println("closing channel")
			return
		}
	}
}

func SendGeneric(out chan<- string, payload Packet, log *logging.Logger) {
	js, err := json.Marshal(payload)
	CheckError(err, log)
	out <- string(js)
}
