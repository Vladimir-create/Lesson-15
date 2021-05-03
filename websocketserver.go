package main

import (
	"net/http"
	"io"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
)

var upgrader = websocket.Upgrader{}
var mas = make([]*websocket.Conn, 0, 0)


func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		
		log.Printf("%s\n", data)
		io.WriteString(w, "successful post")
	} else if req.Method == "OPTIONS" {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(405)
	}
	
}



func Socket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	mas = append(mas, conn)
	if err != nil {
		log.Println(err)
		return
	}
	go writeMess(conn)
	
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("recv: %s", p)
	}
	
}

func writeMess(c *websocket.Conn) {
	var s string
	
	for {
		fmt.Scan(&s)
		err := c.WriteMessage(websocket.TextMessage, []byte(s) )
		if err != nil {
			log.Fatal("dial:", err)
		}
	}
}

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/socket", Socket)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}


