package server

import (
	"context"
	"fmt"

	"github.com/lucheng0127/kube-eip/pkg/manager"
	"github.com/lucheng0127/kube-eip/pkg/protoc/binding"
	ectx "github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
	"github.com/lucheng0127/kube-eip/pkg/utils/validator"
)

const (
	RspSucceed string = "Succeed"
	RspFailed  string = "Failed"
)

type GrpcServer struct {
	binding.UnimplementedEipAgentServer
}

func (s *GrpcServer) EipOperate(ctx context.Context, req *binding.EipOpReq) (*binding.EipOpRsp, error) {
	rsp := new(binding.EipOpRsp)
	rsp.Result = RspSucceed
	rsp.ErrPhase = 0
	tCtx := ectx.NewTraceContext()
	externalIP := validator.ValidateIPv4(req.EipAddr)
	internalIP := validator.ValidateIPv4(req.VmiAddr)
	action := req.Action

	logger.Info(tCtx, fmt.Sprintf("new eip operate request, %s eip %s vmi ip %s",
		action, externalIP.String(), internalIP.String()))

	if externalIP == nil || internalIP == nil {
		logger.Error(tCtx, "invalidate external or internal IP")
		rsp.Result = RspFailed
		rsp.ErrPhase = 0
		return rsp, nil
	}

	mgr := new(manager.EipMgr)
	mgr.ExternalIP = externalIP
	mgr.InternalIP = internalIP
	mgr.IPSetMgr = &manager.IpsetMgr
	mgr.RouteMgr = &manager.RouteMgr

	switch action {
	case "bind":
		errPhase, err := mgr.BindEip()
		if err != nil {
			rsp.Result = RspFailed
			rsp.ErrPhase = int32(errPhase)
			return rsp, nil
		}
	case "unbind":
		errPhase, err := mgr.UnbindEip()
		if err != nil {
			rsp.Result = RspFailed
			rsp.ErrPhase = int32(errPhase)
			return rsp, nil
		}
	default:
		logger.Error(tCtx, "invalidate eip operate")
		rsp.Result = RspFailed
		rsp.ErrPhase = 0
		return rsp, nil
	}

	return rsp, nil
}
