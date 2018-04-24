package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/tominescu/double-golang/simplelog"
)

// thanks to https://github.com/streamlink/streamlink/blob/master/src/streamlink/plugins/douyutv.py
const DOUYU_API_URL_PREFIX = "https://capi.douyucdn.cn/api/v1/"
const DOUYU_API_URL_SUFFIX = "room/%s?aid=wp&cdn=%s&client_sys=wp&time=%d"
const API_SECRET = "zNzMV1y4EMxOHS6I5WKm"

type DouyuResult struct {
	Errno int       `json:"error"`
	Data  DouyuData `json:"data"`
}

type DouyuData struct {
	Status string `json:"show_status"`
	URL    string `json:"hls_url"`
}

func douyuHandler(w http.ResponseWriter, r *http.Request) {
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
	cdns := [...]string{"ws", "tct", "ws2", "dl"}
	now := time.Now().Unix()
	suffix := fmt.Sprintf(DOUYU_API_URL_SUFFIX, id, cdns[0], now)
	hash := md5.New()
	hash.Write([]byte(suffix))
	hash.Write([]byte(API_SECRET))
	sign := fmt.Sprintf("%x", hash.Sum(nil))
	url := DOUYU_API_URL_PREFIX + suffix + "&auth=" + sign
	log.Debug("Douyu api url: %s", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPad; U; CPU OS 3_2_1 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Mobile/7B405")
	client := &http.Client{}
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
	result := DouyuResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
	dst := result.Data.URL
	status := result.Data.Status
	if status != "1" {
		http.Error(w, "房间未开播", 403)
		return
	}
	w.Header().Set("Location", dst)
	http.Error(w, dst, 302)
}
