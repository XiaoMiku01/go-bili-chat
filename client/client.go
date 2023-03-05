package client

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"time"
)
import "github.com/fasthttp/websocket"

type jsonCoder interface {
	Unmarshal(data []byte, v interface{}) error
	Marshal(v interface{}) ([]byte, error)
}

var JsonCoder jsonCoder

type Client struct {
	RoomId    int
	Connected bool
	ctx       context.Context
	cancel    context.CancelFunc
	connect   *websocket.Conn
	revMsg    chan []byte
}

func (c *Client) biliChatConnect(url string) error {
	var err error
	c.connect, _, err = websocket.DefaultDialer.Dial(url, nil)
	if nil != err {
		return err
	}
	return nil
}

func (c *Client) sendAuthMsg(wsAuthMsg WsAuthMessage) error {
	wsPackage := wsAuthMsg.GetPackage()
	err := c.connect.WriteMessage(websocket.BinaryMessage, wsPackage)
	if err != nil {
		return err
	}
	log.Debug("send auth msg to blive success")
	return nil
}

func (c *Client) receiveWsMsg() {
	for {
		select {
		case <-c.ctx.Done():
			log.Debug("receiveWsMsg exit...")
			_ = c.connect.Close()
			return
		default:
			if c.connect != nil && c.Connected {
				_, message, err := c.connect.ReadMessage()
				if err != nil {
					log.Warnf("read blive websocket msg error: %v", err)
					c.Connected = false
					c.connectLoop()
				}
				c.revMsg <- message
			}
		}
	}
}

func (c *Client) heartBeat() {
	for {
		select {
		case <-c.ctx.Done():
			log.Debug("heartBeat exit...")
			_ = c.connect.Close()
			return
		default:
			if c.Connected && c.connect != nil {
				heartBeatPackage := WsHeartBeatMessage{Body: []byte{}}
				log.Debug("send heart beat to blive...")
				_ = c.connect.WriteMessage(websocket.TextMessage, heartBeatPackage.GetPackage())
				time.Sleep(30 * time.Second)
			}
		}
	}
}

func (c *Client) revHandler(handler MsgHandler) {
	for {
		select {
		case <-c.ctx.Done():
			log.Debug("revHandler exit...")
			c.revMsg = nil
			return
		case msg, ok := <-c.revMsg:
			if ok {
				go handler.MsgHandler(msg)
			}
		default:
			time.Sleep(10 * time.Microsecond)
		}
	}
}

func (c *Client) sendConnect() error {
	wsAuthMsg := WsAuthMessage{Body: WsAuthBody{UID: 0, Roomid: c.RoomId, Protover: 3, Platform: "web", Type: 2}}
	u := url.URL{Scheme: "wss", Host: MainWsUrl, Path: "/sub"}
	log.Debug("connect to blive websocket: ", u.String())
	err := c.biliChatConnect(u.String())
	if err != nil {
		apiLiveAuth, err := GetLiveRoomAuth(c.RoomId)
		if err != nil {
			return err
		} else if apiLiveAuth.Code != 0 {
			log.Warnf("get live room info error: %v", apiLiveAuth.Message)
			return RespCodeNotError
		}
		wsAuthMsg.Body.Key = apiLiveAuth.Data.Token
		for nowSum, i := range apiLiveAuth.Data.HostList {
			u := url.URL{Scheme: "wss", Host: i.Host + ":" + strconv.Itoa(i.WssPort), Path: "/sub"}
			log.Debug("connect to blive websocket: ", u.String())
			err = c.biliChatConnect(u.String())
			if err != nil {
				log.Warnf("connect to blive websocket error for %d time: %v\n", nowSum, err)
				if nowSum == 2 {
					return err
				}
			} else {
				log.Debug("connect to blive websocket success")
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	err = c.sendAuthMsg(wsAuthMsg)
	if err != nil {
		log.Warn("send auth msg to websocket error: ", err)
		return err
	}
	return nil
}

func (c *Client) connectLoop() {
	for {
		c.Connected = false
		err := c.sendConnect()
		if err != nil {
			log.Warn("connect to blive error: ", err)
			time.Sleep(5 * time.Second)
		} else {
			log.Info("connected to blive success: ", c.RoomId)
			c.Connected = true
			break
		}
	}
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) BiliChat(CmdChan chan map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("start blive panic: %v", err)
		}
	}()
	c.connectLoop()
	c.revMsg = make(chan []byte, 10)
	handler := MsgHandler{RoomId: c.RoomId, CmdChan: CmdChan}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.revHandler(handler)
	go c.receiveWsMsg()
	go c.heartBeat()
	log.Debug("start blive success", c.RoomId)
}
