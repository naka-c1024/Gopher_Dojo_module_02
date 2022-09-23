package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sync/errgroup"
)

const divNum = 5

func goroutine(url string, arrRange []string) (string, error) {
	var splitData []string = make([]string, divNum)
	eg, ctx := errgroup.WithContext(context.Background())
	for i, ctxRange := range arrRange {
		i := i
		ctxRange := ctxRange
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					return err
				}
				req.Header.Set("Range", ctxRange)
				client := new(http.Client)
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				byteArray, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				splitData[i] = fmt.Sprint(string(byteArray))
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		return "", err
	}
	var allData string
	for _, v := range splitData {
		allData += v
	}
	return allData, nil
}

func hasAcceptRangesBytes(url string) (bool, error) {
	res, err := http.Head(url)
	if err != nil {
		return false, err
	}
	acceptRanges := res.Header["Accept-Ranges"]
	for _, v := range acceptRanges {
		if v == "bytes" {
			return true, nil
		}
	}
	return false, nil
}

func getContentLength(url string) (int, error) {
	res, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	contentLength := res.Header["Content-Length"]
	int64CtLen, err := strconv.ParseInt(contentLength[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(int64CtLen), nil
}

func makeRanges(num int, length int) []string {
	var result []string
	div := length / num
	start := 0
	end := div
	for length > 0 {
		str := fmt.Sprintf("bytes=%d-%d", start, end)
		start = end + 1
		length -= div
		if length < 0 {
			break
		} else {
			end = start + div
		}
		result = append(result, str)
	}
	return result
}

func createFile(url string, content string) error {
	basename := filepath.Base(url)
	f, err := os.Create(basename)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprint(f, content)
	return nil
}

func splitDownload(url string) {
	contentLength, err := getContentLength(url)
	if err != nil {
		log.Fatalln(err)
	}
	arr := makeRanges(divNum, contentLength)
	allData, err := goroutine(url, arr)
	if err != nil {
		log.Fatalln(err)
	}
	err = createFile(url, allData)
	if err != nil {
		log.Fatalln(err)
	}
}

func batchDownload(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = createFile(url, string(byteArray))
	if err != nil {
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
	byteFlag, err := hasAcceptRangesBytes(url)
	if err != nil {
		log.Fatalln(err)
	}
	if byteFlag == true {
		splitDownload(url)
	} else {
		batchDownload(url)
	}
}
