package network

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/qianjunakasumi/Tennouji/internal/pkg/config"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/logger"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqjce"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqtea"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqtlv"

	"go.uber.org/zap"
)

type (
	// serverListReq 服务器列表请求
	serverListReq struct {
		A     byte `jce:"1"`
		B     byte
		C     byte
		D     string
		E     byte
		AppID int32
		IMEI  string
		F     byte
		G     byte
		H     byte
		I     byte
		J     byte
		K     byte
		L     byte
	}

	// serverListRes 服务器列表响应
	serverListRes struct {
		ServerList [][]byte `jce:"2"`
	}

	// server 服务器
	server struct {
		IP   string `jce:"1"`
		Port int64
	}
)

// GetServers 获取服务器
func GetServers() (s []*net.TCPAddr, err error) {

	tea, _ := qqtea.NewCipher([]byte{240, 68, 31, 95, 244, 45, 165, 143, 220, 247, 148, 154, 186, 98, 212, 17})
	req := tea.Encrypt(buildReq())

	resp, err := http.Post(
		"https://configsvr.msf.3g.qq.com/configsvr/serverlist.jsp",
		"",
		bytes.NewReader(req),
	)
	if err != nil {
		logger.Named("网络").Named("获取服务器").Error("请求失败", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	resb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Named("网络").Named("获取服务器").Error("读取失败", zap.Error(err))
		return
	}

	return parseRes(tea.Decrypt(resb)[4:]), nil // [4:] 略过 Length
}

// buildReq 构建请求
func buildReq() []byte {
	return qqtlv.NewWriter(0).Write(qqjce.NewWriter().Write(&qqjce.Packet{
		Version:    2,
		Controller: "ConfigHttp",
		Method:     "HttpServerListReq",
		Data: qqjce.NewWriter().WriteWithDataV2(qqjce.NewWriter().Write(&serverListReq{
			0, 0, 1, "00000", 100,
			config.AppID, config.IMEI,
			0, 0, 0, 0, 0, 0, 1,
		}), "HttpServerListReq", "ConfigHttp.HttpServerListReq").Bytes()}).Bytes(),
	).BytesWithLV()
}

// parseRes 解析响应
func parseRes(jcedata []byte) (srvs []*net.TCPAddr) {

	p := new(qqjce.Packet)
	qqjce.NewReader(jcedata).Read(p)
	data := qqjce.NewReader(p.Data).ReadWithDataV2("HttpServerListRes", "ConfigHttp.HttpServerListRes")

	res := new(serverListRes)
	qqjce.NewReader(data).Read(res)

	for _, v := range res.ServerList {

		s := new(server)
		qqjce.NewReader(qqjce.NewReader(v).ReadStruct()).Read(s)

		if strings.Contains(s.IP, "qq") {
			continue
		}

		srvs = append(srvs, &net.TCPAddr{IP: net.ParseIP(s.IP), Port: int(s.Port)})
	}

	return
}
