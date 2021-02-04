package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/qianjunakasumi/Tennouji/internal/app/protocol/network"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/logger"

	"go.uber.org/zap"
)

// server 服务器
type server struct {
	conn net.Conn // 连接

	servers []*network.Server // 服务器列表

	isOnline bool // 是否在线
}

// connect 连接
func (s *server) connect(re bool) (err error) {

	if !re {
		if s.servers, err = network.GetServers(); err != nil || len(s.servers) == 0 {
			return errors.New("tennouji: 无法获取服务器")
		}
	}

	if err = s.dial(); err != nil {
		return
	}
	go s.listen()

	return
}

// dial 拨号
func (s *server) dial() (err error) {

	for _, srv := range s.servers {

		adr := srv.Name + ":" + strconv.FormatInt(srv.Port, 10)
		s.conn, err = net.Dial("tcp", adr)
		if err != nil {
			logger.Named("网络").Named("拨号").Debug("拨号失败", zap.String("地址", adr))
			continue
		}

		return nil
	}

	return errors.New("tennouji: 所有服务器均不可企及")
}

// listen 监听
func (s *server) listen() {

	for {

		var (
			l      int32
			length = make([]byte, 4)
			_, err = s.conn.Read(length)
		)
		if err != nil {
			logger.Error("服务器已断开连接", zap.Error(err))
			go s.connect(true)
			return
		}
		err = binary.Read(bytes.NewReader(length), binary.BigEndian, &l)
		if err != nil {
			logger.Error("读取LV长度失败", zap.Error(err))
			continue
		}

		pppp := make([]byte, l)
		_, err = s.conn.Read(pppp)
		if err != nil {
			logger.Error("读取失败", zap.Error(err))
			continue
		}

		fmt.Println("收到来包：", string(pppp))

		// TODO 处理数据包

	}

}
