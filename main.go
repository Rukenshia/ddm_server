package main

import (
	"context"
	"log"
	"net"
	"os/exec"

	"github.com/Rukenshia/ddm_server/proto"
	"google.golang.org/grpc"
)

type ddmServer struct{}

func (d *ddmServer) SwitchInput(ctx context.Context, req *proto.SwitchInputRequest) (*proto.SwitchInputResponse, error) {
	log.Println("SwitchInput called")
	c := exec.Command("C:\\Program Files (x86)\\Dell\\Dell Display Manager\\ddm.exe", "/1:SetActiveInput", getInputFromEnum(req.GetInput()))

	if err := c.Run(); err != nil {
		return &proto.SwitchInputResponse{Okay: false, Error: err.Error()}, nil
	}
	return &proto.SwitchInputResponse{Okay: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterDisplayManagerServer(grpcServer, &ddmServer{})
	grpcServer.Serve(lis)
}

func getInputFromEnum(e proto.SwitchInputRequest_InputType) string {
	switch e {
	case proto.SwitchInputRequest_DP1:
		return "DP"
	case proto.SwitchInputRequest_USB_C:
		return "USB-C"
	}

	return ""
}
