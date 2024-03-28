package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lucheng0127/kube-eip/pkg/protoc/admin"
	"github.com/lucheng0127/kube-eip/pkg/protoc/binding"
	"github.com/lucheng0127/kube-eip/pkg/utils/validator"
	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendEipBindingRequest(target, action, eipAddr, vmiAddr string) (*binding.EipOpRsp, error) {
	conn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if !validator.ValidateAction(action) {
		return nil, errors.New("unsupported action")
	}

	if validator.ValidateIPv4(eipAddr) == nil || validator.ValidateIPv4(vmiAddr) == nil {
		return nil, errors.New("invalidate eip or vmi ipv4 address")
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
		return nil, err
	}

	return rsp, nil
}

func Bind(cCtx *cli.Context) error {
	target := cCtx.String("target")
	eipAddr := cCtx.String("eip-ip")
	vmiAddr := cCtx.String("vmi-ip")

	rsp, err := SendEipBindingRequest(target, "bind", eipAddr, vmiAddr)
	if err != nil {
		return err
	}

	fmt.Printf("bind eip: %s vmi ip: %s\nrsp: %+v\n",
		eipAddr, vmiAddr, rsp)
	return nil
}

func Unbind(cCtx *cli.Context) error {
	target := cCtx.String("target")
	eipAddr := cCtx.String("eip-ip")
	vmiAddr := cCtx.String("vmi-ip")

	rsp, err := SendEipBindingRequest(target, "unbind", eipAddr, vmiAddr)
	if err != nil {
		return err
	}

	fmt.Printf("unbind eip: %s vmi ip: %s\nrsp: %+v\n",
		eipAddr, vmiAddr, rsp)
	return nil
}

func TotpQRCode(cCtx *cli.Context) error {
	conn, err := grpc.Dial(cCtx.String("target"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := admin.NewAdminClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rsp, err := client.TotpSec(ctx, &admin.TotpSecReq{})
	if err != nil {
		return err
	}

	qrInfo := fmt.Sprintf("otpauth://totp/kube-eip:admin?secret=%s&issuer=kube-eip&algorithm=SHA1&digits=%d&period=%d", rsp.Secret, rsp.Digits, rsp.Period)
	qrCode, err := qrcode.New(qrInfo, qrcode.Low)
	if err != nil {
		return err
	}

	asciiQr := qrCode.ToSmallString(false)
	fmt.Println(asciiQr)
	return nil
}
