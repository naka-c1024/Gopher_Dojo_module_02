package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"golang.org/x/sync/errgroup"
)

// curl -O https://projects.intra.42.fr/uploads/document/document/10749/19920104_091532.log
// curl -I https://cdn.intra.42.fr/pdf/pdf/51100/ja.subject.pdf
//  -> Accept-Ranges: bytesがあるから分割ダウンロードok
// Content-Length: 1687489 バイトになっている

func goroutine(url string) {
	var eg errgroup.Group
	eg.Go(func() error {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Range", "bytes=100-200")
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", dump)

		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		dumpResp, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", dumpResp)

		// byteArray, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Println(string(byteArray))

		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "" {
		fmt.Fprintf(os.Stderr, "error: empty argument\n")
		os.Exit(1)
	} else if flag.Arg(1) != "" {
		fmt.Fprintf(os.Stderr, "error: multiple arguments\n")
		os.Exit(1)
	}
	url := flag.Arg(0)
	goroutine(url)
}
