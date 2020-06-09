package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	quic "github.com/lucas-clemente/quic-go"

	"github.com/lucas-clemente/quic-go/h2quic"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	multipath := flag.Bool("m", false, "multipath")
	output := flag.String("o", "", "logging output")
	scheduler := flag.String("s", "RR", "scheduler algorithm")
	flag.Parse()
	urls := flag.Args()

	if *verbose {
		utils.SetLogLevel(utils.LogLevelDebug)
	} else {
		utils.SetLogLevel(utils.LogLevelInfo)
	}
	utils.SetLogTimeFormat("")

	if *output != "" {
		logfile, err := os.Create(*output)
		if err != nil {
			panic(err)
		}
		defer logfile.Close()
		log.SetOutput(logfile)
	}

	quicConfig := &quic.Config{
		CreatePaths: *multipath,
		SchedulerAlgorithm: *scheduler,
	}

	utils.Infof("Multipath enabled: %t", *multipath)
	if (*multipath) {
		utils.Infof("MPQUIC scheduler: %s", *scheduler)
	}
	hclient := &http.Client{
		Transport: &h2quic.RoundTripper{QuicConfig: quicConfig},
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, addr := range urls {
		//这里是对每一个url都发起一次get请求
		utils.Infof("GET %s", addr)
		go func(addr string) {
			start := time.Now()
			rsp, err := hclient.Get(addr)
			if err != nil {
				panic(err)
			}
			utils.Infof("Got response for %s\n: %#v", addr, rsp)

			//读取header
			//如果header里面的content-type=application/octet-stream,那就是一个文件，把文件下载下来放在/home/souloh里面
			headerType := rsp.Header.Get("Content-Type")
			if headerType == "application/octet-stream" {
				filename := rsp.Header.Get("Content-Disposition")
				f, err := os.Create("/home/souloh/" + filename)
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(f, rsp.Body)
				if err != nil {
					panic(err)
				}
				elapsed := time.Since(start)
				utils.Infof("Download Completed, Completion Time: \n%s", elapsed)
				wg.Done()
				return
			}

			body := &bytes.Buffer{}
			_, err = io.Copy(body, rsp.Body)
			if err != nil {
				panic(err)
			}
			utils.Infof("Response Body:")
			//这里得到了response body,但是没有作其他的处理,如果我要下载文件的话,就需要根据header处理,
			utils.Infof("%s", body.Bytes())
			elapsed := time.Since(start)
			utils.Infof("Transmission Completed, Completion Time: %s", elapsed)
			wg.Done()
		}(addr)
	} 
	wg.Wait()
}
