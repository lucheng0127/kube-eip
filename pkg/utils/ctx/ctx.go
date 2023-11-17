package ctx

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/lucheng0127/kube-eip/pkg/utils/log"
)

func genUUID() string {
	uuid := "00000000"
	u := make([]byte, 4)
	_, err := rand.Read(u)
	if err == nil {
		uuid = hex.EncodeToString(u)
	}
	return uuid
}

func AddContextTraceID(ctx context.Context) context.Context {
	traceID := genUUID()
	return context.WithValue(ctx, log.MSG_ID, traceID)

}

func NewTraceContext() context.Context {
	ctx := context.Background()
	return AddContextTraceID(ctx)
}
