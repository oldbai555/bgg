package ptuser

import (
	"context"
	"github.com/oldbai555/bgg/internal/_cmd"
	"github.com/oldbai555/bgg/internal/_const"
	"github.com/oldbai555/bgg/internal/bgrpc/srv"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func NewCheckUserInterceptor(cmdList []*_cmd.Cmd) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var inCmd *_cmd.Cmd
		for _, cc := range cmdList {
			if !strings.HasSuffix(strings.ToLower(info.FullMethod), strings.ToLower(strings.TrimPrefix(cc.Path, "/"))) {
				continue
			}
			inCmd = cc
			break
		}

		nCtx := srv.NewGrpcUCtx(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, lbuser.ErrGrpcParseContextFail
		}

		// 设置 traceId
		val := md.Get(_const.GrpcHeaderTraceId)
		if len(val) > 0 {
			nCtx.SetTraceId(val[0])
			log.SetLogHint(val[0])
		}

		// 设置设备 ID
		val = md.Get(_const.GrpcHeaderDeviceId)
		if len(val) > 0 {
			nCtx.SetDeviceId(val[0])
		}

		// 设置 Sid
		val = md.Get(_const.GrpcHeaderSid)
		if len(val) > 0 {
			nCtx.SetSid(val[0])
		}

		// 设置调用人
		val = md.Get(_const.GrpcHeaderAuthType)
		if len(val) > 0 {
			nCtx.SetAuthType(val[0])
		} else {
			// 没传过来 就拿系统生成的
			if inCmd != nil {
				nCtx.SetAuthType(inCmd.GetAuthType())
			} else {
				nCtx.SetAuthType(_cmd.AuthTypeSystem)
			}
		}

		// 不需要校验 直接走
		if nCtx.AuthType() == _cmd.AuthTypeSystem || nCtx.AuthType() == _cmd.AuthTypePublic {
			return handler(nCtx, req)
		}

		userRsp, err := lbuser.GetLoginUser(ctx, &lbuser.GetLoginUserReq{
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
	userSysRsp, err := lbuser.GetLoginUser(ctx, &lbuser.GetLoginUserReq{
		Sid: sid,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return userSysRsp.BaseUser, nil
}
