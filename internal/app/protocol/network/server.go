package network

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/qianjunakasumi/Tennouji/internal/pkg/config"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqjce"
	"github.com/qianjunakasumi/Tennouji/internal/pkg/qqtlv"
	"github.com/qianjunakasumi/qqtea"
)

// serverListRes 服务器列表响应
type serverListRes struct {
	ServerList [][]byte `jce:"2"`
}

// server 服务器
type server struct {
	IP   string `jce:"1"`
	Port int64
}

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
		return
	}
	defer resp.Body.Close()

	resb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	s = parseRes(tea.Decrypt(resb)[4:]) // [4:] 略过 LV
	return
}

// buildReq 构建请求
func buildReq() (req []byte) {

	conf := qqjce.NewWriter(1).
		WriteByte(0).WriteByte(0).WriteByte(1).WriteString("00000").WriteByte(100).
		WriteInt32(config.AppID).WriteString(config.IMEI).
		WriteByte(0).WriteByte(0).WriteByte(0).WriteByte(0).WriteByte(0).WriteByte(0).
		WriteByte(1).BytesWithPack()

	req = qqtlv.NewWriter(0).Write(
		qqjce.NewWriter().Write(qqjce.Packet{
			Version:    2,
			Controller: "ConfigHttp",
			Method:     "HttpServerListReq",
			Data: qqjce.NewWriter(0).WriteMap(qqjce.DataV2{
				"HttpServerListReq": {"ConfigHttp.HttpServerListReq": conf},
			}).Bytes(),
		}).Bytes(),
	).BytesWithLV()

	return
}

// parseRes 解析响应
func parseRes(jcedata []byte) (sers []*net.TCPAddr) {

	packet := new(qqjce.Packet)
	qqjce.NewReader(jcedata).Read(packet)
	dataV2 := qqjce.NewReader(packet.Data).ReadWithDataV2("HttpServerListRes", "ConfigHttp.HttpServerListRes")

	res := new(serverListRes)
	qqjce.NewReader(dataV2).Read(res)

	for _, v := range res.ServerList {

		s := new(server)
		qqjce.NewReader(v[1:]).Read(s)

		if strings.Contains(s.IP, "qq") {
			continue
		}

		sers = append(sers, &net.TCPAddr{
			IP:   net.ParseIP(s.IP),
			Port: int(s.Port),
		})
	}
	return
}
