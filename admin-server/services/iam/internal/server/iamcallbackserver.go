// Package server 实现 pkg/iamcallback.IamCallbackServer，供 chat-rpc/content-rpc 回调查用户
// 展示信息、存量用户枚举、审计日志写入。原来是单体内嵌的一个 zrpc server（和 REST server、
// TaskCallback server 并存），iam-rpc 拆分后整体原样搬到这里，和 iam-rpc 的其他两个 gRPC
// service（Iam、TaskCallback）同一个进程、同一个端口注册（见 services/iam/iam.go），
// 契约不变（和 taskcallbackserver.go 是完全一样的处理方式）。
package server

import (
	"context"
	"database/sql"
	"time"

	"postapocgame/admin-server/services/iam/internal/consts"
	monitoringmodel "postapocgame/admin-server/services/iam/internal/model/monitoring"
	"postapocgame/admin-server/services/iam/internal/repository/registry"
	pb "postapocgame/admin-server/pkg/iamcallback/pb"
)

type IamCallbackServer struct {
	pb.UnimplementedIamCallbackServer
	domain *registry.Domain
}

func NewIamCallbackServer(domain *registry.Domain) *IamCallbackServer {
	return &IamCallbackServer{domain: domain}
}

// FindActiveUserChunk 分批返回未删除、启用中的用户 ID，逻辑与原
// internal/repository/registry/domain.go 里已删除的 iamUserListerAdapter.FindChunk 完全一致
// （原来是 chatdomain.UserLister 的适配实现，chat 域拆分后这段过滤逻辑原样搬到这里）。
func (s *IamCallbackServer) FindActiveUserChunk(ctx context.Context, req *pb.FindActiveUserChunkRequest) (*pb.FindActiveUserChunkResponse, error) {
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
func (s *IamCallbackServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := s.domain.IAM.User.FindByID(ctx, req.UserId)
	if err != nil || user.DeletedAt != 0 {
		return &pb.GetUserProfileResponse{Exists: false}, nil
	}

	resp := &pb.GetUserProfileResponse{
		Exists:    true,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Signature: user.Signature,
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

// RecordAuditLog 写一条审计日志，逻辑迁移自 pkg/audit.RecordAuditLog（该函数原本直接持有
// gateway 的 *svc.ServiceContext，content-rpc/gateway 拆分后拿不到，改成回调这个方法）。
// ip_address/user_agent 是 iam-rpc 拆分时补的字段：content-rpc 的两处调用点原来就传空
// *http.Request，这两个字段留空；gateway 侧自己的 6 个 RBAC 审计调用点有真实的
// *http.Request，用这两个字段填充。
func (s *IamCallbackServer) RecordAuditLog(ctx context.Context, req *pb.RecordAuditLogRequest) (*pb.RecordAuditLogResponse, error) {
	now := time.Now().Unix()
	auditLog := &monitoringmodel.AuditLog{
		UserId:      req.UserId,
		Username:    req.Username,
		AuditType:   req.AuditType,
		AuditObject: req.AuditObject,
		IpAddress:   req.IpAddress,
		UserAgent:   req.UserAgent,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   0,
	}
	if req.DetailJson != "" {
		auditLog.AuditDetail = sql.NullString{String: req.DetailJson, Valid: true}
	}
	if err := s.domain.Monitoring.AuditLog.Create(ctx, auditLog); err != nil {
		return nil, err
	}
	return &pb.RecordAuditLogResponse{}, nil
}

func (s *IamCallbackServer) resolveDepartmentName(ctx context.Context, id uint64) string {
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

func (s *IamCallbackServer) resolveRoleNames(ctx context.Context, roleIDs []uint64) []string {
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
