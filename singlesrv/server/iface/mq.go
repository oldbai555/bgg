/**
 * @Author: zjj
 * @Date: 2024/9/11
 * @Desc:
**/

package iface

import (
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
	"time"
)

type ITopic interface {
	Pub(ctx uctx.IUCtx, obj proto.Message) error
	DeferredPublish(ctx uctx.IUCtx, delay time.Duration, obj proto.Message) error
	Marshal(obj proto.Message) ([]byte, error)
	Unmarshal(buf []byte, obj proto.Message) error
}
