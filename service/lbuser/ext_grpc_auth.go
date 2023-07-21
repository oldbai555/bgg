package lbuser

import (
	"context"
	"github.com/oldbai555/bgg/pkg/_const"
	"github.com/oldbai555/bgg/pkg/cmd"
	"github.com/oldbai555/lbtool/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CheckUser() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		nCtx := NewUContext(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrGrpcParseContextFail
		}

		// 设置 traceId
		val := md.Get(_const.HeaderLBTraceId)
		if len(val) > 0 {
			nCtx.SetTraceId(val[0])
			log.SetLogHint(val[0])
		}

		// 设置设备 ID
		val = md.Get(_const.HeaderLBDeviceId)
		if len(val) > 0 {
			nCtx.SetDeviceId(val[0])
		}

		// 设置 Sid
		val = md.Get(_const.HeaderLBSid)
		if len(val) > 0 {
			nCtx.SetSid(val[0])
		}

		// 设置调用人
		val = md.Get(_const.HeaderLBCallFrom)
		if len(val) > 0 {
			nCtx.SetCallFrom(val[0])
		} else {
			nCtx.SetCallFrom(cmd.AuthTypeSystem) // 默认就是系统调用
		}

		// 不需要校验 直接走
		if nCtx.CallFrom() == cmd.AuthTypeSystem || nCtx.CallFrom() == cmd.AuthTypePublic {
			// 前置校验
			return handler(nCtx, req)
		}

		// 拿用户信息
		userRsp, err := GetLoginUser(ctx, &GetLoginUserReq{
			Sid: nCtx.Sid(),
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		nCtx.SetUser(userRsp.BaseUser)

		log.Infof("user context is %v", nCtx)

		// 前置校验
		return handler(nCtx, req)
	}
}

type UContext struct {
	context.Context

	sid      string
	deviceId string
	traceId  string
	callFrom string

	user *BaseUser
}

func (U *UContext) CallFrom() string {
	return U.callFrom
}

func (U *UContext) SetCallFrom(callFrom string) {
	U.callFrom = callFrom
}

func (U *UContext) User() *BaseUser {
	return U.user
}

func (U *UContext) SetUser(user *BaseUser) {
	U.user = user
}

func NewUContext(ctx context.Context) *UContext {
	return &UContext{
		Context: ctx,
	}
}

func ToUCtx(ctx context.Context) *UContext {
	return ctx.(*UContext)
}

func (U *UContext) Sid() string {
	return U.sid
}

func (U *UContext) SetSid(sid string) {
	U.sid = sid
}

func (U *UContext) DeviceId() string {
	return U.deviceId
}

func (U *UContext) SetDeviceId(deviceId string) {
	U.deviceId = deviceId
}

func (U *UContext) TraceId() string {
	return U.traceId
}

func (U *UContext) SetTraceId(traceId string) {
	U.traceId = traceId
}
