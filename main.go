package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	requestsCount, _ := strconv.Atoi(os.Args[1])
	maxConcurrency, _ := strconv.Atoi(os.Args[2])
	url := os.Args[3]

	client := CreateClient()
	request, _ := CreateRequest(url)

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
			go DoRequest(client, request, concurrencyMonitor)
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

func CreateClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: (&net.Dialer{
			Timeout: 60 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   120 * time.Second,
		ResponseHeaderTimeout: 180 * time.Second,
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   0,
		DisableKeepAlives:     true,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   200 * time.Second,
	}

	return client
}

func CreateRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request: %w", err)
	}
	request.Header.Add("Accept", "application/json")

	return request, nil
}

func DoRequest(client *http.Client, request *http.Request, monitor chan string) {
	response, err := client.Do(request)
	if err != nil {
		monitor <- err.Error()
		return
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		monitor <- err.Error()
		return
	}

	_ = response.Body.Close()
	monitor <- "OK"
}
