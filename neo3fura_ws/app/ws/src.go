package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"neo3fura_ws/cli"
	"neo3fura_ws/home"
	"net/http"

	"github.com/gorilla/websocket"
)

var add = flag.String("addr", "0.0.0.0:2026", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func mainpage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DETECT CONNECTION")
	client := &cli.T{}
	c := &home.T{
		Client: client,
	}
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	mt, _, err := wsc.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	var responseChannel = make(chan map[string]interface{}, 20)

	go c.GetAddressCount(&responseChannel)
	go c.GetAssetCount(&responseChannel)
	go c.GetBlockCount(&responseChannel)
	go c.GetBlockInfoList(&responseChannel)
	go c.GetCandidateCount(&responseChannel)
	go c.GetContractCount(&responseChannel)
	go c.GetTransactionCount(&responseChannel)
	go c.GetTransactionList(&responseChannel)
	go ResponseController(mt, wsc, &responseChannel)
}

func ResponseController(mt int, wsc *websocket.Conn, ch *chan map[string]interface{}) {
	str := "hello websocket"
	err := wsc.WriteMessage(mt, []byte(str))
	if err != nil {
		log.Fatal(err)
	}
	for {
		b := <-*ch
		sent, err := json.Marshal(b)
		if err != nil {
			log.Fatal(err)
		}
		err = wsc.WriteMessage(mt, sent)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/home", mainpage)
	log.Fatal(http.ListenAndServe(*add, nil))
}
