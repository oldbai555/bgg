package initdata

// 种子数据固定 ID 常量（对应 db/services/iam/**/init_*.sql 里幂等写入的初始化记录）。
// 这些 ID 是脚手架/初始化 SQL 约定死的常量，不是运行时可配置值，因此就地定义在本包，
// 不依赖任何 services/<name>/internal/consts（那些是各服务的 Go internal 包，本包在
// services/ 之外，按 Go 可见性规则本来就无法导入）。
const (
	// superAdminUserID 超级管理员账号的固定 ID
	superAdminUserID uint64 = 1
	// superAdminRoleID 超级管理员角色的固定 ID
	superAdminRoleID uint64 = 1
	// rootDepartmentID 根部门的固定 ID
	rootDepartmentID uint64 = 1

	// 初始化权限ID：1, 10-13, 20-23, 30-33, 40-43, 50-53, 60-63
	permissionIDAll         uint64 = 1
	permissionIDRangeStart1 uint64 = 10
	permissionIDRangeEnd1   uint64 = 13
	permissionIDRangeStart2 uint64 = 20
	permissionIDRangeEnd2   uint64 = 23
	permissionIDRangeStart3 uint64 = 30
	permissionIDRangeEnd3   uint64 = 33
	permissionIDRangeStart4 uint64 = 40
	permissionIDRangeEnd4   uint64 = 43
	permissionIDRangeStart5 uint64 = 50
	permissionIDRangeEnd5   uint64 = 53
	permissionIDRangeStart6 uint64 = 60
	permissionIDRangeEnd6   uint64 = 63

	// 初始化菜单ID：1, 10-16
	menuIDRoot       uint64 = 1
	menuIDRangeStart uint64 = 10
	menuIDRangeEnd   uint64 = 16
)

// IsInitUserID 检查是否是初始化用户ID
func IsInitUserID(id uint64) bool {
	return id == superAdminUserID
}

// IsInitRoleID 检查是否是初始化角色ID
func IsInitRoleID(id uint64) bool {
	return id == superAdminRoleID
}

// IsInitPermissionID 检查是否是初始化权限ID
func IsInitPermissionID(id uint64) bool {
	if id == permissionIDAll {
		return true
	}
	if id >= permissionIDRangeStart1 && id <= permissionIDRangeEnd1 {
		return true
	}
	if id >= permissionIDRangeStart2 && id <= permissionIDRangeEnd2 {
		return true
	}
	if id >= permissionIDRangeStart3 && id <= permissionIDRangeEnd3 {
		return true
	}
	if id >= permissionIDRangeStart4 && id <= permissionIDRangeEnd4 {
		return true
	}
	if id >= permissionIDRangeStart5 && id <= permissionIDRangeEnd5 {
		return true
	}
	if id >= permissionIDRangeStart6 && id <= permissionIDRangeEnd6 {
		return true
	}
	return false
}

// IsInitDepartmentID 检查是否是初始化部门ID
func IsInitDepartmentID(id uint64) bool {
	return id == rootDepartmentID
}

// IsInitMenuID 检查是否是初始化菜单ID
func IsInitMenuID(id uint64) bool {
	return id == menuIDRoot || (id >= menuIDRangeStart && id <= menuIDRangeEnd)
}
