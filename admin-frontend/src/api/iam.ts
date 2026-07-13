import {
  // 认证
  login,
  refresh,
  logout,
  profile,
  profileUpdate,
  passwordChange,
  // 用户
  userList,
  userCreate,
  userUpdate,
  userDelete,
  userRoleList,
  userRoleUpdate,
  // 角色
  roleList,
  roleCreate,
  roleUpdate,
  roleDelete,
  rolePermissionList,
  rolePermissionUpdate,
  // 权限
  permissionList,
  permissionCreate,
  permissionUpdate,
  permissionDelete,
  permissionMenuList,
  permissionMenuUpdate,
  permissionApiList,
  permissionApiUpdate,
  // 部门
  departmentTree,
  departmentCreate,
  departmentUpdate,
  departmentDelete,
  // 菜单
  menuTree,
  menuMyTree,
  menuCreate,
  menuUpdate,
  menuDelete,
  // API 管理
  apiList,
  apiCreate,
  apiUpdate,
  apiDelete
} from '@/api/generated/admin'

/**
 * IAM 域 API 封装（用户/角色/权限/部门/菜单/API 管理）
 */
export const iamApi = {
  // ========== 认证 ==========
  login,
  refresh,
  logout,
  profile,
  profileUpdate,
  passwordChange,

  // ========== 用户 ==========
  userList,
  userCreate,
  userUpdate,
  userDelete,
  userRoleList,
  userRoleUpdate,

  // ========== 角色 ==========
  roleList,
  roleCreate,
  roleUpdate,
  roleDelete,
  rolePermissionList,
  rolePermissionUpdate,

  // ========== 权限 ==========
  permissionList,
  permissionCreate,
  permissionUpdate,
  permissionDelete,
  permissionMenuList,
  permissionMenuUpdate,
  permissionApiList,
  permissionApiUpdate,

  // ========== 部门 ==========
  departmentTree,
  departmentCreate,
  departmentUpdate,
  departmentDelete,

  // ========== 菜单 ==========
  menuTree,
  menuMyTree,
  menuCreate,
  menuUpdate,
  menuDelete,

  // ========== API 管理 ==========
  apiList,
  apiCreate,
  apiUpdate,
  apiDelete
}
