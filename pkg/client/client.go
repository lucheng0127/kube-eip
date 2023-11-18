package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lucheng0127/kube-eip/pkg/protoc/binding"
	"github.com/lucheng0127/kube-eip/pkg/utils/validator"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Launch(cCtx *cli.Context) error {
	// Dial server
	conn, err := grpc.Dial(cCtx.String("target"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	action := cCtx.String("action")
	eipAddr := cCtx.String("eip-ip")
	vmiAddr := cCtx.String("vmi-ip")

	if !validator.ValidateAction(action) {
		return errors.New("unsupported action")
	}

	if validator.ValidateIPv4(eipAddr) == nil || validator.ValidateIPv4(vmiAddr) == nil {
		return errors.New("invalidate eip or vmi ipv4 address")
	}

	// Build client
	client := binding.NewEipAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send request
	rsp, err := client.EipOperate(ctx, &binding.EipOpReq{
		Action:  action,
		EipAddr: eipAddr,
		VmiAddr: vmiAddr,
	})

	if err != nil {
		return err
	}

	fmt.Printf("action: %s eip: %s vmi ip: %s\nrsp: %+v\n",
		action, eipAddr, vmiAddr, rsp)
	return nil
}
