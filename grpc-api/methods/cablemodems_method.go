package methods

import (
	"context"
	"strconv"

	"api-project/grpc-api/gen/cablemodems"
)

// CableModemMethod 实现 CableModemServiceServer 接口
type CableModemMethod struct {
	cablemodems.UnimplementedCableModemServiceServer
}

// 示例实现：GetCableModem
func (h *CableModemMethod) ByMac(ctx context.Context, req *cablemodems.ByMacRequest) (*cablemodems.ByMacResponse, error) {
	// 暂时返回 mock 数据
	res := &cablemodems.ByMacResponse{
		Modems: make([]*cablemodems.CableModem, 0, len(req.MacAddress)),
	}
	for index, cmMac := range req.MacAddress {
		res.Modems = append(res.Modems, &cablemodems.CableModem{
			Mac:   cmMac,
			Model: "ExampleModel" + strconv.Itoa(index),
			Ipv4:  "192.168.100.1",
		})
	}
	return res, nil
}
