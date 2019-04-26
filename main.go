package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/Rukenshia/ddm_server/proto"
	"github.com/getlantern/systray"
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
	onExit := func() {
		os.Exit(0)
	}
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("DDM Server")
	systray.SetTooltip("Dell Display Manager Server")
	quit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-quit.ClickedCh:
				systray.Quit()
			}
		}
	}()

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
