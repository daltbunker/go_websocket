package main

import (
  "golang.org/x/net/websocket"
  "io"
  "fmt"
  "log"
  "net/http"
)

type Server struct {
 conn map[*websocket.Conn]bool
}


func newServer() *Server {
  return &Server {
    conn: make(map[*websocket.Conn]bool), 
  }
}

func (s *Server) handleWS(ws *websocket.Conn) {
  fmt.Println("New connection from client: ", ws.RemoteAddr()) 

  // map is not safe for concurrency, replace map with Mutex
  s.conn[ws] = true
  s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
  buf := make([]byte, 1024)
  for {
    n, err := ws.Read(buf) 
    if err != nil {
      // case when client closes connection??
      if err == io.EOF {
        break
      }
      fmt.Println("Read error: ", err)
      // if return here, connection will be dropped 
      continue
    }
    msg := buf[:n]
    fmt.Println(string(msg))

    ws.Write([]byte("Thanks for the message!"))
  } 
}

func main() {
  server := newServer()
  http.Handle("/ws", websocket.Handler(server.handleWS))
  log.Fatal(http.ListenAndServe(":8080", nil))
}

