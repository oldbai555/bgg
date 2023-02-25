package impl

import (
	"context"
	"github.com/oldbai555/bgg/client/lbuser"
)

var lbuserServer LbuserServer

type LbuserServer struct {
	*lbuser.UnimplementedLbuserServer
}

func (a *LbuserServer) ResetPassword(ctx context.Context, req *lbuser.ResetPasswordReq) (*lbuser.ResetPasswordRsp, error) {
	var rsp lbuser.ResetPasswordRsp
	return &rsp, nil
}
func (a *LbuserServer) Logout(ctx context.Context, req *lbuser.LogoutReq) (*lbuser.LogoutRsp, error) {
	var rsp lbuser.LogoutRsp
	return &rsp, nil
}
func (a *LbuserServer) GetLoginUser(ctx context.Context, req *lbuser.GetLoginUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp
	return &rsp, nil
}
func (a *LbuserServer) GetUserList(ctx context.Context, req *lbuser.GetUserListReq) (*lbuser.GetUserListRsp, error) {
	var rsp lbuser.GetUserListRsp
	return &rsp, nil
}
func (a *LbuserServer) GetUser(ctx context.Context, req *lbuser.GetUserReq) (*lbuser.GetUserRsp, error) {
	var rsp lbuser.GetUserRsp
	return &rsp, nil
}
func (a *LbuserServer) UpdateUserNameWithRole(ctx context.Context, req *lbuser.UpdateUserNameWithRoleReq) (*lbuser.UpdateUserNameWithRoleRsp, error) {
	var rsp lbuser.UpdateUserNameWithRoleRsp
	return &rsp, nil
}
func (a *LbuserServer) Login(ctx context.Context, req *lbuser.LoginReq) (*lbuser.LoginRsp, error) {
	var rsp lbuser.LoginRsp
	return &rsp, nil
}
func (a *LbuserServer) UpdateLoginUserInfo(ctx context.Context, req *lbuser.UpdateLoginUserInfoReq) (*lbuser.UpdateLoginUserInfoRsp, error) {
	var rsp lbuser.UpdateLoginUserInfoRsp
	return &rsp, nil
}
func (a *LbuserServer) AddUser(ctx context.Context, req *lbuser.AddUserReq) (*lbuser.AddUserRsp, error) {
	var rsp lbuser.AddUserRsp
	return &rsp, nil
}
func (a *LbuserServer) DelUser(ctx context.Context, req *lbuser.DelUserReq) (*lbuser.DelUserRsp, error) {
	var rsp lbuser.DelUserRsp
	return &rsp, nil
}
func (a *LbuserServer) GetFrontUser(ctx context.Context, req *lbuser.GetFrontUserReq) (*lbuser.GetLoginUserRsp, error) {
	var rsp lbuser.GetLoginUserRsp
	return &rsp, nil
}
