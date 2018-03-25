package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const MY_COPY_BUF_SIZE = 1024 * 4

var gclient = &http.Client{}

func MyCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, MY_COPY_BUF_SIZE)
	for {
		n, err := src.Read(buf)
		written += int64(n)
		nn := 0
		for nn < n {
			nw, ew := dst.Write(buf[nn:n])
			nn += nw
			if ew != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func MultiDownload(w io.Writer, url string, threadNum int) (written int64, err error) {
	res, err := http.Head(url)
	if err != nil {
		return
	}
	length, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return
	}
	q := length / threadNum
	r := length % threadNum
	doneChan := make(chan int, threadNum)
	body := make([]string, threadNum)
	done := make([]int, threadNum)
	curr := -1
	for i := 0; i < threadNum; i++ {
		min := i * q
		max := (i + 1) * q
		if i == threadNum-1 {
			max += r
		}
		go Download(doneChan, body, url, i, min, max)
	}

	for i := range doneChan {
		if i == -1 {
			return written, fmt.Errorf("Download partly failed")
		}
		done[i] = 1
		for {
			if done[curr+1] == 0 {
				break
			}
			curr++
			n, _ := w.Write([]byte(body[curr]))
			written += int64(n)
			if curr == threadNum-1 {
				return
			}
		}
	}
	return
}

func Download(done chan int, body []string, url string, i, min, max int) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		done <- -1
		return err
	}
	range_header := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
	req.Header.Add("Range", range_header)
	resp, err := gclient.Do(req)
	if err != nil {
		done <- -1
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		done <- -1
		return err
	}
	body[i] = string(content)
	done <- i
	return nil
}