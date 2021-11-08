package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {

	requestsCount, _ := strconv.Atoi(os.Args[1])
	maxConcurrency, _ := strconv.Atoi(os.Args[2])
	host := os.Args[3]

	dialer := &net.Dialer{
		Timeout: 60 * time.Second,
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	concurrencyMonitor := make(chan string, maxConcurrency)

	concurrent := 0

	requestsCountOk := 0
	requestsCountFailed := 0
	requestsLastError := ""

	lastMonitoringMessageTime := time.Now().Add(-3 * time.Second)

	for {
		if time.Now().Sub(lastMonitoringMessageTime) > 2*time.Second {
			fmt.Printf("remaining requests: %v, concurrent requests: %v, ok: %v, failed: %v, lastError: %v\n", requestsCount, concurrent, requestsCountOk, requestsCountFailed, requestsLastError)
			lastMonitoringMessageTime = time.Now()
			requestsCountOk = 0
			requestsCountFailed = 0
			requestsLastError = ""
		}

		if concurrent < maxConcurrency && requestsCount > 0 {
			concurrent++
			requestsCount--
			go DoRequest(dialer, config, host, concurrencyMonitor)
		} else {
			time.Sleep(10 * time.Microsecond)
		}

		select {
		case status := <-concurrencyMonitor:
			concurrent--
			if status == "OK" {
				requestsCountOk++
			} else {
				requestsCountFailed++
				requestsLastError = status
			}
		default:
		}

		if requestsCount <= 0 && concurrent <= 0 {
			break
		}

	}
}

func DoRequest(dialer *net.Dialer, config *tls.Config, host string, monitor chan string) {
	conn, err := tls.DialWithDialer(dialer, "tcp", host, config)
	if err != nil {
		monitor <- err.Error()
		return
	}

	_ = conn.Close()
	monitor <- "OK"
}
