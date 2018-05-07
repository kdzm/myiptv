package main

import (
	"io"
	"net/http"

	log "github.com/tominescu/double-golang/simplelog"
)

func ttcatvHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request URL:%s", r.URL)
	url := "http://httpdvb.slave.ttcatv.tv:13164" + r.URL.String()
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
	}
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)

	/*
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, http.StatusText(503), 503)
			return
		}
		body = bytes.Replace(body, []byte("httpdvb.slave.ttcatv.tv:13164"), []byte(r.Host), -1)
		w.Write(body)
	*/
}
