package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type tpingResult struct {
	time time.Duration
	err  error
}

type tpingStat struct {
	connResult  *tpingResult
	closeResult *tpingResult
}

func tping(host string, port int, timeout time.Duration) (*tpingResult, *tpingResult) {
	connResult := &tpingResult{}
	closeResult := &tpingResult{}
	// var result strings.Builder

	t1 := time.Now()
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	diff := time.Now().Sub(t1)
	if err != nil {
		connResult.err = err
		return connResult, nil
	}
	connResult.time = diff

	t1 = time.Now()
	err = conn.Close()
	diff = time.Now().Sub(t1)
	if err != nil {
		closeResult.err = err
		return connResult, closeResult
	}
	closeResult.time = diff
	return connResult, closeResult
}

func tpingString(connResult *tpingResult, closeResult *tpingResult) string {
	var result strings.Builder
	if connResult.err != nil {
		result.WriteString(fmt.Sprintf("[%s] \tconnection error: (%s)", connResult.time, connResult.err))
		return result.String()
	}
	result.WriteString(fmt.Sprintf("[%s] \tconnect ok\t-\t", connResult.time))

	if closeResult.err != nil {
		result.WriteString(fmt.Sprintf("close error: %s \t[%s]", closeResult.err, closeResult.time))
		return result.String()
	}
	result.WriteString(fmt.Sprintf("close ok \t[%s]", closeResult.time))
	return result.String()
}
