package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"api-project/grpc-api/gen/cablemodems"

	"google.golang.org/grpc"
)

type CableModemsClient struct {
	client cablemodems.CableModemServiceClient
}

func NewCableModemsClient(conn *grpc.ClientConn) *CableModemsClient {
	return &CableModemsClient{
		client: cablemodems.NewCableModemServiceClient(conn),
	}
}

func (c *CableModemsClient) ByMac(macs []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.ByMac(ctx, &cablemodems.ByMacRequest{
		MacAddress: macs,
	})
	if err != nil {
		log.Fatalf("could not get cable modem: %v", err)
	}
	test, _ := json.Marshal(resp)
	log.Printf("Response: ", string(test))
}
