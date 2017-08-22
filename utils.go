package bchainlibs

import (
	"math/rand"
	"time"
	"strconv"
	"github.com/op/go-logging"
	"os"
	"fmt"
)

var src = rand.NewSource(time.Now().UnixNano())

var LogFormat = logging.MustStringFormatter(
	"%{level:.4s}=> %{time:0102 15:04:05.999} %{shortfile} %{message}",
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

func WaitForSync( targetSync float64, log *logging.Logger ) {
	// ------------

	// It gives some time for the network to get configured before it gets its own IP.
	// This value should be passed as a environment variable indicating the time when
	// the simulation starts, this should be calculated by an external source so all
	// Go programs containers start at the same UnixTime.
	now := float64(time.Now().Unix())
	sleepTime := 0
	if targetSync > now {
		sleepTime = int(targetSync - now)
		log.Info("SYNC: Sync time is " + strconv.FormatFloat( targetSync, 'f', 6, 64) )
	}
	//else {
	//sleepTime = globalNumberNodes
	//}
	log.Info("SYNC: sleepTime is " + strconv.Itoa(sleepTime))
	time.Sleep(time.Second * time.Duration(sleepTime))
	// ------------
}

func PrepareLog( logName string ) (*os.File) {
	var logPath = "/var/log/golang/"
	//var logPath = "./golang/"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, 0777)
	}

	var logFile = logPath + logName + ".log"
	f, err := os.OpenFile(logFile, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	return f
}