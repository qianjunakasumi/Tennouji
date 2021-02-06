package network

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/qianjunakasumi/Tennouji/internal/app/config"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/logger"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqjce"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqtea"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqtlv"

	"go.uber.org/zap"
)

type (
	// serverListReq 服务器列表请求
	serverListReq struct {
		Number  int64  `jce:"1"` // 号码
		Timeout byte   // 超时时间
		C       byte   // 未知字段
		IMSI    string // 国际移动用户识别码
		ISWIFI  bool   // 是否 WIFI 环境
		AppID   int32  // App ID
		IMEI    string // 国际移动设备识别码
		CellID  byte   // 基站编号？
	}

	// serverListRes 服务器列表响应
	serverListRes struct {
		ServerListWithWIFI [][]byte `jce:"3"` // WIFI模式下服务器列表
	}

	// Server 服务器
	Server struct {
		Host string `jce:"1"` // 主机
		Port int64  // 端口
		City string `jce:"8"` // 城市
	}
)

// GetServers 获取服务器
func GetServers() (s []*Server, err error) {

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
	return qqtlv.NewWriter(0).WriteBytes(qqjce.NewWriter().Write(&qqjce.Packet{
		Version:    3,
		Controller: "ConfigHttp",
		Method:     "HttpServerListReq",
		Data: qqjce.NewWriter().WriteWithDataV3(
			qqjce.NewWriter().Write(&serverListReq{ // TODO 支持更多的参数
				0, 0, 1, "", true,
				config.AppID, config.IMEI, 0,
			}),
			"HttpServerListReq",
		).Bytes(),
	}).Bytes()).BytesWithLV()
}

// parseRes 解析响应
func parseRes(jcedata []byte) (srvs []*Server) {

	p := new(qqjce.Packet)
	qqjce.NewReader(jcedata).Read(p)
	data := qqjce.NewReader(p.Data).ReadWithDataV3("HttpServerListRes")

	res := new(serverListRes)
	qqjce.NewReader(data).Read(res)

	for _, v := range res.ServerListWithWIFI {
		s := new(Server)
		qqjce.NewReader(qqjce.NewReader(v).ReadStruct()).Read(s)

		srvs = append(srvs, s)
	}

	return
}
