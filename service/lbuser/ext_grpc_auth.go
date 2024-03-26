package lbuser

import (
	"context"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/bcmd"
	"github.com/oldbai555/micro/bconst"
	"github.com/oldbai555/micro/brpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func NewCheckUserInterceptor(cmdList []*bcmd.Cmd) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var inCmd *bcmd.Cmd
		for _, cc := range cmdList {
			if !strings.HasSuffix(strings.ToLower(info.FullMethod), strings.ToLower(strings.TrimPrefix(cc.Path, "/"))) {
				continue
			}
			inCmd = cc
			break
		}

		nCtx := brpc.NewGrpcUCtx(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrGrpcParseContextFail
		}

		// 设置 traceId
		val := md.Get(bconst.GrpcHeaderTraceId)
		if len(val) > 0 {
			nCtx.SetTraceId(val[0])
			log.SetLogHint(val[0])
		}

		// 设置设备 ID
		val = md.Get(bconst.GrpcHeaderDeviceId)
		if len(val) > 0 {
			nCtx.SetDeviceId(val[0])
		}

		// 设置 Sid
		val = md.Get(bconst.GrpcHeaderSid)
		if len(val) > 0 {
			nCtx.SetSid(val[0])
		}

		// 设置调用人
		val = md.Get(bconst.GrpcHeaderAuthType)
		if len(val) > 0 {
			nCtx.SetAuthType(val[0])
		} else {
			// 没传过来 就拿系统生成的
			if inCmd != nil {
				nCtx.SetAuthType(inCmd.GetAuthType())
			} else {
				nCtx.SetAuthType(bcmd.AuthTypeSystem)
			}
		}

		// 不需要校验 直接走
		if nCtx.AuthType() == bcmd.AuthTypeSystem || nCtx.AuthType() == bcmd.AuthTypePublic {
			return handler(nCtx, req)
		}

		userRsp, err := GetLoginUser(ctx, &GetLoginUserReq{
			Sid: nCtx.Sid(),
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		nCtx.SetExtInfo(userRsp.BaseUser)

		log.Infof("user context is %v", nCtx)

		// 前置校验
		return handler(nCtx, req)
	}
}

func CheckLoginUser(ctx context.Context, sid string) (interface{}, error) {
	userSysRsp, err := GetLoginUser(ctx, &GetLoginUserReq{
		Sid: sid,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return userSysRsp.BaseUser, nil
}
