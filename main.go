package main

import (
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gotping"
	app.Usage = "TCP connect/close ping"
	app.ArgsUsage = "host port"
	app.Version = "0.1.0"

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	app.Flags = []cli.Flag{
		cli.Float64Flag{Name: "timeout, t", Value: 1},
		cli.Float64Flag{Name: "sleep, s", Value: 1},
		// cli.IntFlag{Name: "parallel, p", Value: 1},
		cli.IntFlag{Name: "count, c", Value: 0},
	}
	app.Action = func(c *cli.Context) error {
		host := c.Args().Get(0)
		port, err := strconv.Atoi(c.Args().Get(1))
		if err != nil {
			return err
		}

		timeoutSeconds := c.Float64("timeout")
		timeout := time.Duration(timeoutSeconds) * time.Second
		sleepSeconds := c.Float64("sleep")
		sleep := time.Duration(sleepSeconds) * time.Second
		// parallel := c.Int("parallel")
		count := c.Int("count")
		countStr := "infinite"
		if count != 0 {
			countStr = fmt.Sprintf("%d", count)
		}

		fmt.Printf("TCP Ping to %s:%d\nCount: %s\nTimeout: %.3f second(s)\nSleep: %.3f second(s)\n\n",
			host, port, countStr, timeoutSeconds, sleepSeconds)

		incrBy := 1
		if count == 0 {
			count = 1
			incrBy = 0
		}
		tpingStatList := list.New()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		go func() {
			<-ch
			printStats(host, port, tpingStatList)
			os.Exit(0)
		}()

		for i := 0; i < count; i += incrBy {
			connResult, closeResult := tping(host, port, timeout)
			tpingStatList.PushBack(tpingStat{connResult, closeResult})
			fmt.Printf("%s\n", tpingString(connResult, closeResult))
			if incrBy == 0 || i != count-1 {
				time.Sleep(sleep)
			}
		}

		printStats(host, port, tpingStatList)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printStats(host string, port int, tpingStatList *list.List) {
	fmt.Printf("\n--- %s:%d statistics ---\n\n", host, port)

	var connTotal, closeTotal, connOk, closeOk int
	var connMinDuration, connMaxDuration, connTotalDuration time.Duration
	var closeMinDuration, closeMaxDuration, closeTotalDuration time.Duration
	connMinDuration = math.MaxInt64
	closeMinDuration = math.MaxInt64

	for e := tpingStatList.Front(); e != nil; e = e.Next() {
		s := e.Value.(tpingStat)

		if s.connResult != nil {
			connTotal++
			if s.connResult.err == nil {
				connTotalDuration += s.connResult.time
				if s.connResult.time < connMinDuration {
					connMinDuration = s.connResult.time
				}
				if s.connResult.time > connMaxDuration {
					connMaxDuration = s.connResult.time
				}
				connOk++
			}
		}

		if s.closeResult != nil {
			closeTotal++
			if s.closeResult.err == nil {
				closeTotalDuration += s.closeResult.time
				if s.closeResult.time < closeMinDuration {
					closeMinDuration = s.closeResult.time
				}
				if s.closeResult.time > closeMaxDuration {
					closeMaxDuration = s.closeResult.time
				}
				closeOk++
			}
		}
	}
	fmt.Printf("# CONNECT:\ttotal:%d, failed:%d (%.1f%%), min:%s, max:%s, average:%s\n",
		connTotal, connTotal-connOk, 100*float64(connTotal-connOk)/float64(connTotal),
		connMinDuration, connMaxDuration, connTotalDuration/time.Duration(connOk))
	fmt.Printf("# CLOSE:\ttotal:%d, failed:%d (%.1f%%), min:%s, max:%s, average:%s\n",
		closeTotal, closeTotal-closeOk, 100*float64(closeTotal-closeOk)/float64(closeTotal),
		closeMinDuration, closeMaxDuration, closeTotalDuration/time.Duration(closeOk))
}
