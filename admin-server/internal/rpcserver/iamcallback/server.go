// Package iamcallback 实现 pkg/iamcallback.IamCallbackServer，供已经拆分出去、但还没有
// iam-rpc 可以直接调用的服务（当前是 chat-rpc）回调查用户数据。当前阶段单体内嵌一个 zrpc
// server 提前实现这份接口（admin.go 里和 REST server、TaskCallback server 并存），后续
// iam-rpc 真正拆分时把这个实现原样搬过去，不改契约（和 internal/rpcserver/taskcallback
// 是完全一样的处理方式）。
package iamcallback

import (
	"context"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/repository/registry"
	pb "postapocgame/admin-server/pkg/iamcallback/pb"
)

type Server struct {
	pb.UnimplementedIamCallbackServer
	domain *registry.Domain
}

func NewServer(domain *registry.Domain) *Server {
	return &Server{domain: domain}
}

// FindActiveUserChunk 分批返回未删除、启用中的用户 ID，逻辑与原
// internal/repository/registry/domain.go 里已删除的 iamUserListerAdapter.FindChunk 完全一致
// （原来是 chatdomain.UserLister 的适配实现，chat 域拆分后这段过滤逻辑原样搬到这里）。
func (s *Server) FindActiveUserChunk(ctx context.Context, req *pb.FindActiveUserChunkRequest) (*pb.FindActiveUserChunkResponse, error) {
	users, newLastID, err := s.domain.IAM.User.FindChunk(ctx, req.Limit, req.LastId)
	if err != nil {
		return nil, err
	}

	refs := make([]*pb.ActiveUserRef, 0, len(users))
	for _, u := range users {
		if u.DeletedAt != 0 || u.Status != consts.UserStatusEnabled {
			continue
		}
		refs = append(refs, &pb.ActiveUserRef{Id: u.Id})
	}
	return &pb.FindActiveUserChunkResponse{Users: refs, NextLastId: newLastID}, nil
}

// GetUserProfile 返回 chat-rpc 展示用的用户信息（用户名/昵称/头像/部门名/角色名列表），
// 逻辑迁移自原 internal/logic/chat/{chat,group}/*.go 里重复出现的
// "FindByID + Department.ListAll 建 map + Role.FindPage 建 map + ListRoleIDsByUserID" 组合。
func (s *Server) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := s.domain.IAM.User.FindByID(ctx, req.UserId)
	if err != nil || user.DeletedAt != 0 {
		return &pb.GetUserProfileResponse{Exists: false}, nil
	}

	resp := &pb.GetUserProfileResponse{
		Exists:   true,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}

	if user.DepartmentId > 0 {
		resp.DepartmentName = s.resolveDepartmentName(ctx, user.DepartmentId)
	}

	roleIDs, err := s.domain.IAM.UserRole.ListRoleIDsByUserID(ctx, user.Id)
	if err == nil && len(roleIDs) > 0 {
		resp.RoleNames = s.resolveRoleNames(ctx, roleIDs)
	}

	return resp, nil
}

func (s *Server) resolveDepartmentName(ctx context.Context, id uint64) string {
	depts, err := s.domain.IAM.Department.ListAll(ctx)
	if err != nil {
		return ""
	}
	for _, d := range depts {
		if d.Id == id && d.DeletedAt == 0 {
			return d.Name
		}
	}
	return ""
}

func (s *Server) resolveRoleNames(ctx context.Context, roleIDs []uint64) []string {
	// FindPage(1, 10000, "") 拿全量角色列表建 id->name 映射，和原 chat logic 里的写法一致
	// （角色总数量级很小，10000 上限沿用既有代码的做法，不是本次新引入的性能问题）。
	roles, _, err := s.domain.IAM.Role.FindPage(ctx, 1, 10000, "")
	if err != nil {
		return nil
	}
	roleMap := make(map[uint64]string, len(roles))
	for _, r := range roles {
		if r.DeletedAt == 0 {
			roleMap[r.Id] = r.Name
		}
	}
	names := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		if name, ok := roleMap[id]; ok {
			names = append(names, name)
		}
	}
	return names
}
