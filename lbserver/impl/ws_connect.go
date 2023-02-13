package impl

import (
	"github.com/gorilla/websocket"
	"github.com/oldbai555/bgg/lbwebsocket"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"net/http"
	"sync"
)

var wu = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		//允许跨域
		return true
	},
}

// Connection ws连接
type Connection struct {
	// websocket连接
	wsConn *websocket.Conn
	// 读取websocket的channel
	inChan chan []byte
	// 给websocket写消息的channel
	outChan chan []byte
	// 关闭连接管道信号
	closeChan chan byte
	// 锁
	mutex sync.Mutex
	// closeChan 状态
	isClosed bool
}

// InitConnection 初始化长连接
func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	//启动读协程
	go conn.readLoop()
	//启动写协程
	go conn.writeLoop()
	return
}

// ReadMessage 读取websocket消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = lberr.NewErr(int32(lbwebsocket.ErrCode_ErrCodeConnectClosed), "connection is closed")
	}
	return
}

// WriteMessage 发送消息到websocket
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = lberr.NewErr(int32(lbwebsocket.ErrCode_ErrCodeConnectClosed), "connection is closed")
	}
	return
}

// Close 关闭连接
func (conn *Connection) Close() {
	//线程安全的Close,可重入
	err := conn.wsConn.Close()
	if err != nil {
		log.Errorf("err is : %v", err)
		return
	}

	//只执行一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 等待读数据
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		// 如果数据量过大阻塞在这里,等待inChan有空闲的位置！
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			//closeChan关闭的时候
			goto ERR
		}
	}
ERR:
	conn.Close()
}

// 等待写数据
func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
