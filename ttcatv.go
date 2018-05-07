package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/tominescu/double-golang/simplelog"
)

func ttcatvHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request URL:%s", r.URL)
	url := "http://httpdvb.slave.ttcatv.tv:13164" + r.URL.String()
	var resp *http.Response
	var err error
	for i := 0; i < 10; i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	defer resp.Body.Close()

	if strings.HasSuffix(r.URL.Path, ".ts") {
		io.Copy(w, resp.Body)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	body = bytes.Replace(body, []byte("httpdvb.slave.ttcatv.tv:13164"), []byte(r.Host), -1)
	w.Write(body)
}
