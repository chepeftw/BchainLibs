package bchainlibs

import (
	"math/rand"
	"time"
	"strconv"
	"github.com/op/go-logging"
	"os"
	"fmt"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"net"
)

var src = rand.NewSource(time.Now().UnixNano())

var LogFormat = logging.MustStringFormatter(
	"%{level:.4s}=> %{time:0102 15:04:05.999} %{shortfile} %{message}",
)

var LogFormatPimp = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func MyCalculateSHA(str string) string {
	t := sha256.Sum256([]byte( str ))
	t2 := sha256.Sum256(t[:])
	checksum := hex.EncodeToString(t2[:])

	return checksum
}

func WaitForSync(targetSync float64, log *logging.Logger) {
	// ------------

	// It gives some time for the network to get configured before it gets its own IP.
	// This value should be passed as a environment variable indicating the time when
	// the simulation starts, this should be calculated by an external source so all
	// Go programs containers start at the same UnixTime.
	now := float64(time.Now().Unix())
	sleepTime := 0
	if targetSync > now {
		sleepTime = int(targetSync - now)
		log.Info("SYNC: Sync time is " + strconv.FormatFloat(targetSync, 'f', 6, 64))
	}
	//else {
	//sleepTime = globalNumberNodes
	//}
	log.Info("SYNC: sleepTime is " + strconv.Itoa(sleepTime))
	time.Sleep(time.Second * time.Duration(sleepTime))
	// ------------
}

func PrepareLogGen(logConfPath string, logName string, extension string) (*os.File) {
	var logPath = logConfPath
	if logConfPath == "" {
		logPath = "/var/log/golang/"
	}

	if !strings.HasSuffix(logPath, "/") {
		logPath = logPath + "/"
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, 0777)
	}

	var logFile = logPath + logName + "." + extension
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	return f
}

func PrepareLog(logConfPath string, logName string) (*os.File) {
	return PrepareLogGen( logConfPath, logName, "log")
}

// A Simple function to verify error
func CheckError(err error, log *logging.Logger) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

// Getting my own IP, first we get all interfaces, then we iterate
// discard the loopback and get the IPv4 address, which should be the eth0
func SelfieIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}

	return net.ParseIP("127.0.0.1")
}

func CompareIPs(a net.IP, b net.IP) bool {
	if a == nil || b == nil {
		return false
	}

	return a.String() == b.String()
}
