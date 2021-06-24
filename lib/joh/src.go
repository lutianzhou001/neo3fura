package joh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"neo3fura/lib/rwio"
	"neo3fura/lib/scex"
	"net/http"
	"net/http/httputil"
	"net/rpc"
	"net/url"
	"path/filepath"
	"sort"
	// "sort"
)

// T ...
type T struct {
}

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
		log.Printf("Error in reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
	}
	r := req.Clone(req.Context())
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Error decoding in JOSN: %v", err)
		http.Error(w, "can't decoding in JSON", http.StatusBadRequest)
	}

	c, err := me.OpenConfigFile()
	if err != nil {
		log.Fatalln(err)
	}

	sort.Strings(c.Methods.Realized)
	index := sort.SearchStrings(c.Methods.Realized, fmt.Sprintf("%v", request["method"]))
	if index < len(c.Methods.Realized) && c.Methods.Realized[index] == request["method"] {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		conn := &rwio.T{R: req.Body, W: w}
		codec := &scex.T{}
		codec.Init(conn)
		rpc.ServeCodec(codec)
	} else {
		me.Handle(c.Proxy.URI[repostMode], w, r)
		repostMode = (repostMode + 1) % 5
	}
}

func (me *T) Handle(target string, w http.ResponseWriter, r *http.Request) {
	uri, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(uri)
	r.URL.Host = uri.Host
	r.URL.Scheme = uri.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = uri.Host
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
