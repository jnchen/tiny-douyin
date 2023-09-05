package service

import (
	"douyin/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var chatConnMap = sync.Map{}

func RunMessageServer(quit <-chan struct{}) {
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Printf("Run message sever failed: %v\n", err)
		return
	}

	go func() {
		for {
			conn, err := listen.Accept()
			if errors.Is(err, net.ErrClosed) {
				break
			} else if err != nil {
				log.Printf("Accept conn failed: %v\n", err)
				continue
			}

			go process(conn)
		}
	}()

	<-quit
	log.Println("Message server received quit signal. Shutting down...")
	if err = listen.Close(); err != nil {
		log.Printf("Error while closing message server: %v\n", err)
	}
}

func process(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	var buf [256]byte
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				break
			}
			log.Printf("Read message failed: %v\n", err)
			continue
		}

		var event = model.MessageSendEvent{}
		_ = json.Unmarshal(buf[:n], &event)
		log.Printf("Receive Message: %+v\n", event)

		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
		if len(event.MsgContent) == 0 {
			chatConnMap.Store(fromChatKey, conn)
			continue
		}

		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
		writeConn, exist := chatConnMap.Load(toChatKey)
		if !exist {
			log.Printf("User %d offline\n", event.ToUserId)
			continue
		}

		pushEvent := model.MessagePushEvent{
			FromUserId: event.UserId,
			MsgContent: event.MsgContent,
		}
		pushData, _ := json.Marshal(pushEvent)
		_, err = writeConn.(net.Conn).Write(pushData)
		if err != nil {
			log.Printf("Push message failed: %v\n", err)
		}
	}
}
