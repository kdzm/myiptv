package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/tominescu/double-golang/simplelog"
)

const TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJkZXZpY2VJZCI6IjI0MTQ3Y2JmLTA0ZTAtMzE4MS1iZjZhLWI5OWNhZjNkZjdjNSIsInRpbWVzdGFtcCI6MTUyOTU5NjA4MiwiY2hhbm5lbElkIjoiNTRhMjZkNDEtYTBkMi00MDVhLWIwMWEtMjgwZjU0OWQ4YjFlIn0.tlOkaViZk3SSC-I5Uin_7jTY9nDkx9-gDn2cSR7BJDM"

type BestvResult struct {
	Data []BestvData `json:"data"`
}

type BestvData struct {
	Live string `json:"live"`
}

func bestvHandler(w http.ResponseWriter, r *http.Request) {
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

	url := "https://bestvapi.bestv.cn/video/live_rate?tid=" + id + "&se=weixin&ct=3&d=3&_fk=0&token=" + TOKEN
	resp, err := http.Get(url)
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
	result := BestvResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

	var live string
	for _, data := range result.Data {
		live = data.Live
	}
	if live == "" {
		http.Error(w, "no url found", 503)
		return
	}
	w.Header().Set("Location", live)
	http.Error(w, live, 302)
}
