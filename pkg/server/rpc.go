package server

import (
	"context"

	"github.com/lucheng0127/kube-eip/pkg/protoc/binding"
)

type GrpcServer struct {
	binding.UnimplementedEipAgentServer
}

func (GrpcServer) mustEmbedUnimplementedEipAgentServer() {}

func (s *GrpcServer) EipOperate(ctx context.Context, req *binding.EipOpReq) (*binding.EipOpRsp, error) {
	rsp := new(binding.EipOpRsp)
	rsp.Result = "Succeed"
	rsp.ErrPhase = 0
	return rsp, nil
}
