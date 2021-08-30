package joh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"neo3fura_http/config"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/rwio"
	"neo3fura_http/lib/scex"
	"net/http"
	"net/http/httputil"
	"net/rpc"
	"net/url"
	"path/filepath"
	"sort"
	// "sort"
)

// T ...
type T struct{}

type Config struct {
	Methods struct {
		Realized []string `yaml:"realized"`
	} `yaml:"methods"`
	Proxy struct {
		URI []string `yaml:"uri"`
	} `yaml:"proxy"`
}

// To repost to every nodes in queue
var repostMode int = 0

func (me *T) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log2.Infof("Error in reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
	}
	r := req.Clone(req.Context())
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {
		log2.Infof("Error decoding in JSON: %v", err)
		http.Error(w, "can't decoding in JSON", http.StatusBadRequest)
	} else {
		log2.Infof("Request is: %v", request["method"])
		c, err := me.OpenConfigFile()
		if err != nil {
			log2.Fatalf("open config file error:%s", err)
		}
		sort.Strings(config.Apis)
		index := sort.SearchStrings(config.Apis, fmt.Sprintf("%v", request["method"]))
		if index < len(config.Apis) && config.Apis[index] == request["method"] {
			// can find
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			conn := &rwio.T{R: req.Body, W: w}
			codec := &scex.T{}
			codec.Init(conn)
			rpc.ServeCodec(codec)
		} else {
			// can't find
			me.Handle(c.Proxy.URI[repostMode], w, r)
			repostMode = (repostMode + 1) % 5
		}
	}
}

func (me *T) Handle(target string, w http.ResponseWriter, r *http.Request) {
	log2.Infof("Repost to node")
	uri, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(uri)
	r.URL.Host = uri.Host
	r.URL.Scheme = uri.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = uri.Host
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	proxy.ServeHTTP(w, r)
}

func (me *T) OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("./config.yml")
	f, err := ioutil.ReadFile(absPath)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}
