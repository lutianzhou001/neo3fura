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
	"net/rpc"
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
		log2.Infof("index is %v", index)
		log2.Infof("lens of apis is %v", len(config.Apis))
		log2.Infof("config api is %v", config.Apis[index])
		log2.Infof("request[method] is %v", request["method"])
		log2.Infof("index < len(config.apis) %v", index < len(config.Apis))
		log2.Infof("config.apis == request[method] %v", config.Apis[index] == request["method"])
		if index < len(config.Apis) && config.Apis[index] == request["method"] {
			// can find
			log2.Infof("Serving %v", request["method"])
			conn := &rwio.T{R: req.Body, W: w}
			codec := &scex.T{}
			codec.Init(conn)
			rpc.ServeCodec(codec)
		} else {
			// can't find
			log2.Infof("repost %v", request["method"])
			responseBody := bytes.NewBuffer(body)
			w.Header().Set("Content-Type", "application/json")
			resp, err := http.Post(c.Proxy.URI[repostMode], "application/json", responseBody)
			if err != nil {
				log2.Fatalf("Repost error%v", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log2.Fatalf("Read err%v", err)
			}
			w.Write(body)
			repostMode = (repostMode + 1) % 5
		}
	}
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
