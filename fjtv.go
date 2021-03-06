package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/tominescu/double-golang/simplelog"
)

func fjtvApiHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request URL:%s", r.URL)
	url := "http://stream6.fjtv.net" + r.URL.String()
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func fjtvHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request URL:%s", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	id := r.Form.Get("id")
	if id == "" {
		http.Error(w, http.StatusText(503), 503)
		return
	}
	uri := "http://live.fjtv.net/m2o/channel/channel_info.php?id=" + id
	client := &http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "curl/7.52.1")
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	re := regexp.MustCompile(`http:[^\"]*\.m3u8\?_upt=\w*`)
	hls := re.Find(body)
	dst := strings.Replace(string(hls), "\\", "", -1)
	req, _ = http.NewRequest("GET", dst, nil)
	req.Header.Set("User-Agent", "curl/7.52.1")
	resp2, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	defer resp2.Body.Close()
	body, err = ioutil.ReadAll(resp2.Body)
	re = regexp.MustCompile(`.*\.m3u8\?_upt=.*`)
	hls = re.Find(body)
	u, err := url.Parse(string(hls))
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	base, err := url.Parse(dst)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	dst = strings.Replace(base.ResolveReference(u).String(), "stream6.fjtv.net", r.Host, -1)
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
