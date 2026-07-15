import webapi from "./gocliRequest"
import * as components from "./adminComponents"
export * from "./adminComponents"

/**
 * @description 
 * @param req
 */
export function blogArticleList(req: components.BlogArticleListReq) {
	return webapi.get<components.BlogArticleListResp>(`/api/v1/blog/articles`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleCreate(req: components.BlogArticleCreateReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleUpdate(req: components.BlogArticleUpdateReq) {
	return webapi.put<components.Response>(`/api/v1/blog/articles`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleDelete(req: components.BlogArticleDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/blog/articles`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleDetail(req: components.BlogArticleDetailReq) {
	return webapi.get<components.BlogArticleDetailResp>(`/api/v1/blog/articles/detail`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticlePublish(req: components.BlogArticlePublishReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/publish`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleSubmit(req: components.BlogArticleSubmitReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/submit`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleTop(req: components.BlogArticleTopReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/top`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleUnpublish(req: components.BlogArticleUnpublishReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/unpublish`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleUntop(req: components.BlogArticleUntopReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/untop`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleAudit(req: components.BlogArticleAuditReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/audit`, req)
}

/**
 * @description 
 * @param req
 */
export function blogArticleAuditUnpublish(req: components.BlogArticleAuditUnpublishReq) {
	return webapi.post<components.Response>(`/api/v1/blog/articles/audit/unpublish`, req)
}

/**
 * @description 
 * @param req
 */
export function blogFriendLinkList(req: components.BlogFriendLinkListReq) {
	return webapi.get<components.BlogFriendLinkListResp>(`/api/v1/blog/friend-links`, req)
}

/**
 * @description 
 * @param req
 */
export function blogFriendLinkCreate(req: components.BlogFriendLinkCreateReq) {
	return webapi.post<components.Response>(`/api/v1/blog/friend-links`, req)
}

/**
 * @description 
 * @param req
 */
export function blogFriendLinkUpdate(req: components.BlogFriendLinkUpdateReq) {
	return webapi.put<components.Response>(`/api/v1/blog/friend-links`, req)
}

/**
 * @description 
 * @param req
 */
export function blogFriendLinkDelete(req: components.BlogFriendLinkDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/blog/friend-links`, req)
}

/**
 * @description 
 */
export function publicBlogArticleStats() {
	return webapi.get<components.PublicBlogArticleStatsResp>(`/api/v1/public/blog/article-stats`)
}

/**
 * @description 
 * @param req
 */
export function publicBlogArticleList(req: components.PublicBlogArticleListReq) {
	return webapi.get<components.PublicBlogArticleListResp>(`/api/v1/public/blog/articles`, req)
}

/**
 * @description 
 * @param req
 */
export function publicBlogArticleDetail(req: components.PublicBlogArticleDetailReq) {
	return webapi.get<components.PublicBlogArticleDetailResp>(`/api/v1/public/blog/articles/info`, req)
}

/**
 * @description 
 * @param req
 */
export function publicBlogArticleNext(req: components.PublicBlogArticleNextReq) {
	return webapi.get<components.PublicBlogArticleNextResp>(`/api/v1/public/blog/articles/next`, req)
}

/**
 * @description 
 * @param req
 */
export function publicBlogArticlePrev(req: components.PublicBlogArticlePrevReq) {
	return webapi.get<components.PublicBlogArticlePrevResp>(`/api/v1/public/blog/articles/prev`, req)
}

/**
 * @description 
 */
export function publicBlogAuthorInfo() {
	return webapi.get<components.PublicBlogAuthorInfoResp>(`/api/v1/public/blog/author-info`)
}

/**
 * @description 
 */
export function publicBlogFriendLinkList() {
	return webapi.get<components.PublicBlogFriendLinkListResp>(`/api/v1/public/blog/friend-links`)
}

/**
 * @description 
 */
export function publicBlogSocialInfoList() {
	return webapi.get<components.PublicBlogSocialInfoListResp>(`/api/v1/public/blog/social-infos`)
}

/**
 * @description 
 */
export function publicBlogTagList() {
	return webapi.get<components.PublicBlogTagListResp>(`/api/v1/public/blog/tags`)
}

/**
 * @description 
 * @param req
 */
export function blogSocialInfoList(req: components.BlogSocialInfoListReq) {
	return webapi.get<components.BlogSocialInfoListResp>(`/api/v1/blog/social-infos`, req)
}

/**
 * @description 
 * @param req
 */
export function blogSocialInfoCreate(req: components.BlogSocialInfoCreateReq) {
	return webapi.post<components.Response>(`/api/v1/blog/social-infos`, req)
}

/**
 * @description 
 * @param req
 */
export function blogSocialInfoUpdate(req: components.BlogSocialInfoUpdateReq) {
	return webapi.put<components.Response>(`/api/v1/blog/social-infos`, req)
}

/**
 * @description 
 * @param req
 */
export function blogSocialInfoDelete(req: components.BlogSocialInfoDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/blog/social-infos`, req)
}

/**
 * @description 
 * @param req
 */
export function blogTagList(req: components.BlogTagListReq) {
	return webapi.get<components.BlogTagListResp>(`/api/v1/blog/tags`, req)
}

/**
 * @description 
 * @param req
 */
export function blogTagCreate(req: components.BlogTagCreateReq) {
	return webapi.post<null>(`/api/v1/blog/tags`, req)
}

/**
 * @description 
 * @param req
 */
export function blogTagUpdate(req: components.BlogTagUpdateReq) {
	return webapi.put<null>(`/api/v1/blog/tags`, req)
}

/**
 * @description 
 * @param req
 */
export function blogTagDelete(req: components.BlogTagDeleteReq) {
	return webapi.delete<null>(`/api/v1/blog/tags`, req)
}

/**
 * @description 
 * @param req
 */
export function blogTagOptions(req: components.BlogTagOptionsReq) {
	return webapi.get<components.BlogTagOptionsResp>(`/api/v1/blog/tags/options`, req)
}

/**
 * @description 
 */
export function chatList() {
	return webapi.get<components.ChatListResp>(`/api/v1/chats`)
}

/**
 * @description 
 * @param req
 */
export function chatMessageSend(req: components.ChatMessageSendReq) {
	return webapi.post<components.ChatMessageSendResp>(`/api/v1/chats/messages`, req)
}

/**
 * @description 
 * @param req
 */
export function chatMessageList(req: components.ChatMessageListReq) {
	return webapi.get<components.ChatMessageListResp>(`/api/v1/chats/messages/list`, req)
}

/**
 * @description 
 * @param params
 */
export function chatGroupList(params: components.ChatGroupListReqParams) {
	return webapi.get<components.ChatGroupListResp>(`/api/v1/chats/groups`, params)
}

/**
 * @description 
 * @param req
 */
export function chatGroupCreate(req: components.ChatGroupCreateReq) {
	return webapi.post<components.Response>(`/api/v1/chats/groups`, req)
}

/**
 * @description 
 * @param req
 */
export function chatGroupUpdate(req: components.ChatGroupUpdateReq) {
	return webapi.put<components.Response>(`/api/v1/chats/groups`, req)
}

/**
 * @description 
 * @param req
 */
export function chatGroupDelete(req: components.ChatGroupDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/chats/groups`, req)
}

/**
 * @description 
 * @param req
 */
export function chatGroupDetail(req: components.ChatGroupDetailReq) {
	return webapi.get<components.ChatGroupDetailResp>(`/api/v1/chats/groups/detail`, req)
}

/**
 * @description 
 * @param req
 */
export function chatGroupMemberList(req: components.ChatGroupMemberListReq) {
	return webapi.get<components.ChatGroupMemberListResp>(`/api/v1/chats/groups/members`, req)
}

/**
 * @description 
 * @param req
 */
export function chatGroupMemberAdd(req: components.ChatGroupMemberAddReq) {
	return webapi.post<components.Response>(`/api/v1/chats/groups/members`, req)
}

/**
 * @description 
 * @param req
 */
export function chatGroupMemberRemove(req: components.ChatGroupMemberRemoveReq) {
	return webapi.delete<components.Response>(`/api/v1/chats/groups/members`, req)
}

/**
 * @description 
 * @param req
 */
export function chatMessageListAdmin(req: components.ChatMessageListReq) {
	return webapi.get<components.ChatMessageListResp>(`/api/v1/chats/messages`, req)
}

/**
 * @description 
 * @param req
 */
export function chatMessageDelete(req: components.ChatMessageDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/chats/messages`, req)
}

/**
 * @description 
 * @param req
 */
export function apiList(req: components.ApiListReq) {
	return webapi.get<components.ApiListResp>(`/api/v1/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function apiCreate(req: components.ApiCreateReq) {
	return webapi.post<null>(`/api/v1/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function apiUpdate(req: components.ApiUpdateReq) {
	return webapi.put<null>(`/api/v1/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function apiDelete(req: components.ApiDeleteReq) {
	return webapi.delete<null>(`/api/v1/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function login(req: components.LoginReq) {
	return webapi.post<components.TokenPair>(`/api/v1/login`, req)
}

/**
 * @description 
 * @param req
 */
export function refresh(req: components.RefreshReq) {
	return webapi.post<components.TokenPair>(`/api/v1/refresh`, req)
}

/**
 * @description 
 * @param req
 */
export function logout(req: components.LogoutReq) {
	return webapi.post<null>(`/api/v1/logout`, req)
}

/**
 * @description 
 */
export function profile() {
	return webapi.get<components.ProfileResp>(`/api/v1/profile`)
}

/**
 * @description 
 * @param req
 */
export function profileUpdate(req: components.ProfileUpdateReq) {
	return webapi.put<null>(`/api/v1/profile`, req)
}

/**
 * @description 
 * @param req
 */
export function passwordChange(req: components.PasswordChangeReq) {
	return webapi.post<null>(`/api/v1/profile/password`, req)
}

/**
 * @description 
 * @param req
 */
export function departmentCreate(req: components.DepartmentCreateReq) {
	return webapi.post<null>(`/api/v1/departments`, req)
}

/**
 * @description 
 * @param req
 */
export function departmentUpdate(req: components.DepartmentUpdateReq) {
	return webapi.put<null>(`/api/v1/departments`, req)
}

/**
 * @description 
 * @param req
 */
export function departmentDelete(req: components.DepartmentDeleteReq) {
	return webapi.delete<null>(`/api/v1/departments`, req)
}

/**
 * @description 
 */
export function departmentTree() {
	return webapi.get<components.DepartmentTreeResp>(`/api/v1/departments/tree`)
}

/**
 * @description 
 * @param req
 */
export function menuCreate(req: components.MenuCreateReq) {
	return webapi.post<null>(`/api/v1/menus`, req)
}

/**
 * @description 
 * @param req
 */
export function menuUpdate(req: components.MenuUpdateReq) {
	return webapi.put<null>(`/api/v1/menus`, req)
}

/**
 * @description 
 * @param req
 */
export function menuDelete(req: components.MenuDeleteReq) {
	return webapi.delete<null>(`/api/v1/menus`, req)
}

/**
 * @description 
 */
export function menuTree() {
	return webapi.get<components.MenuTreeResp>(`/api/v1/menus/tree`)
}

/**
 * @description 
 */
export function menuMyTree() {
	return webapi.get<components.MenuTreeResp>(`/api/v1/menus/my-tree`)
}

/**
 * @description 
 * @param req
 */
export function permissionList(req: components.PermissionListReq) {
	return webapi.get<components.PermissionListResp>(`/api/v1/permissions`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionCreate(req: components.PermissionCreateReq) {
	return webapi.post<null>(`/api/v1/permissions`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionUpdate(req: components.PermissionUpdateReq) {
	return webapi.put<null>(`/api/v1/permissions`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionDelete(req: components.PermissionDeleteReq) {
	return webapi.delete<null>(`/api/v1/permissions`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionApiList(req: components.PermissionApiListReq) {
	return webapi.get<components.PermissionApiListResp>(`/api/v1/permissions/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionApiUpdate(req: components.PermissionApiUpdateReq) {
	return webapi.put<null>(`/api/v1/permissions/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionMenuList(req: components.PermissionMenuListReq) {
	return webapi.get<components.PermissionMenuListResp>(`/api/v1/permissions/menus`, req)
}

/**
 * @description 
 * @param req
 */
export function permissionMenuUpdate(req: components.PermissionMenuUpdateReq) {
	return webapi.put<null>(`/api/v1/permissions/menus`, req)
}

/**
 * @description 
 * @param req
 */
export function roleList(req: components.RoleListReq) {
	return webapi.get<components.RoleListResp>(`/api/v1/roles`, req)
}

/**
 * @description 
 * @param req
 */
export function roleCreate(req: components.RoleCreateReq) {
	return webapi.post<null>(`/api/v1/roles`, req)
}

/**
 * @description 
 * @param req
 */
export function roleUpdate(req: components.RoleUpdateReq) {
	return webapi.put<null>(`/api/v1/roles`, req)
}

/**
 * @description 
 * @param req
 */
export function roleDelete(req: components.RoleDeleteReq) {
	return webapi.delete<null>(`/api/v1/roles`, req)
}

/**
 * @description 
 * @param req
 */
export function rolePermissionList(req: components.RolePermissionListReq) {
	return webapi.get<components.RolePermissionListResp>(`/api/v1/roles/permissions`, req)
}

/**
 * @description 
 * @param req
 */
export function rolePermissionUpdate(req: components.RolePermissionUpdateReq) {
	return webapi.put<null>(`/api/v1/roles/permissions`, req)
}

/**
 * @description 
 * @param req
 */
export function userList(req: components.UserListReq) {
	return webapi.get<components.UserListResp>(`/api/v1/users`, req)
}

/**
 * @description 
 * @param req
 */
export function userCreate(req: components.UserCreateReq) {
	return webapi.post<null>(`/api/v1/users`, req)
}

/**
 * @description 
 * @param req
 */
export function userUpdate(req: components.UserUpdateReq) {
	return webapi.put<null>(`/api/v1/users`, req)
}

/**
 * @description 
 * @param req
 */
export function userDelete(req: components.UserDeleteReq) {
	return webapi.delete<null>(`/api/v1/users`, req)
}

/**
 * @description 
 * @param req
 */
export function userRoleList(req: components.UserRoleListReq) {
	return webapi.get<components.UserRoleListResp>(`/api/v1/users/roles`, req)
}

/**
 * @description 
 * @param req
 */
export function userRoleUpdate(req: components.UserRoleUpdateReq) {
	return webapi.put<null>(`/api/v1/users/roles`, req)
}

/**
 * @description 
 * @param req
 */
export function dailyShortSentenceList(req: components.DailyShortSentenceListReq) {
	return webapi.get<components.DailyShortSentenceListResp>(`/api/v1/daily-short-sentences`, req)
}

/**
 * @description 
 * @param req
 */
export function dailyShortSentenceCreate(req: components.DailyShortSentenceCreateReq) {
	return webapi.post<null>(`/api/v1/daily-short-sentences`, req)
}

/**
 * @description 
 * @param req
 */
export function dailyShortSentenceUpdate(req: components.DailyShortSentenceUpdateReq) {
	return webapi.put<null>(`/api/v1/daily-short-sentences`, req)
}

/**
 * @description 
 * @param req
 */
export function dailyShortSentenceDelete(req: components.DailyShortSentenceDeleteReq) {
	return webapi.delete<null>(`/api/v1/daily-short-sentences`, req)
}

/**
 * @description 
 * @param req
 */
export function demoList(req: components.DemoListReq) {
	return webapi.get<components.DemoListResp>(`/api/v1/demos`, req)
}

/**
 * @description 
 * @param req
 */
export function demoCreate(req: components.DemoCreateReq) {
	return webapi.post<null>(`/api/v1/demos`, req)
}

/**
 * @description 
 * @param req
 */
export function demoUpdate(req: components.DemoUpdateReq) {
	return webapi.put<null>(`/api/v1/demos`, req)
}

/**
 * @description 
 * @param req
 */
export function demoDelete(req: components.DemoDeleteReq) {
	return webapi.delete<null>(`/api/v1/demos`, req)
}

/**
 * @description 
 */
export function ping() {
	return webapi.get<components.PingResp>(`/api/v1/ping`)
}

/**
 * @description 
 * @param req
 */
export function publicDictGet(req: components.DictGetReq) {
	return webapi.get<components.DictGetResp>(`/api/v1/public/dict`, req)
}

/**
 * @description 
 * @param req
 */
export function auditLogList(req: components.AuditLogListReq) {
	return webapi.get<components.AuditLogListResp>(`/api/v1/audit-logs`, req)
}

/**
 * @description 
 * @param req
 */
export function auditLogDetail(req: components.AuditLogDetailReq) {
	return webapi.get<components.AuditLogDetailResp>(`/api/v1/audit-logs/detail`, req)
}

/**
 * @description 
 * @param req
 */
export function auditLogExport(req: components.AuditLogExportReq) {
	return webapi.get<components.AuditLogExportResp>(`/api/v1/audit-logs/export`, req)
}

/**
 * @description 
 * @param req
 */
export function loginLogList(req: components.LoginLogListReq) {
	return webapi.get<components.LoginLogListResp>(`/api/v1/login-logs`, req)
}

/**
 * @description 
 * @param req
 */
export function loginLogDetail(req: components.LoginLogDetailReq) {
	return webapi.get<components.LoginLogDetailResp>(`/api/v1/login-logs/detail`, req)
}

/**
 * @description 
 * @param req
 */
export function loginLogExport(req: components.LoginLogExportReq) {
	return webapi.get<components.LoginLogExportResp>(`/api/v1/login-logs/export`, req)
}

/**
 * @description 
 */
export function loginLogStats() {
	return webapi.get<components.LoginLogStatsResp>(`/api/v1/login-logs/stats`)
}

/**
 * @description 
 * @param req
 */
export function metricReport(req: components.MetricReportReq) {
	return webapi.post<components.Response>(`/api/v1/metrics/report`, req)
}

/**
 * @description 
 */
export function metricReportOptions() {
	return webapi.options<null>(`/api/v1/metrics/report`)
}

/**
 * @description 
 * @param req
 */
export function metricStats(req: components.MetricStatsReq) {
	return webapi.get<components.MetricStatsResp>(`/api/v1/metrics/stats`, req)
}

/**
 * @description 
 */
export function monitorStats() {
	return webapi.get<components.MonitorStatsResp>(`/api/v1/monitor/stats`)
}

/**
 * @description 
 */
export function monitorStatus() {
	return webapi.get<components.MonitorStatusResp>(`/api/v1/monitor/status`)
}

/**
 * @description 
 * @param req
 */
export function operationLogList(req: components.OperationLogListReq) {
	return webapi.get<components.OperationLogListResp>(`/api/v1/operation-logs`, req)
}

/**
 * @description 
 * @param req
 */
export function operationLogDetail(req: components.OperationLogDetailReq) {
	return webapi.get<components.OperationLogDetailResp>(`/api/v1/operation-logs/detail`, req)
}

/**
 * @description 
 * @param req
 */
export function operationLogExport(req: components.OperationLogExportReq) {
	return webapi.get<components.OperationLogExportResp>(`/api/v1/operation-logs/export`, req)
}

/**
 * @description 
 * @param req
 */
export function performanceLogList(req: components.PerformanceLogListReq) {
	return webapi.get<components.PerformanceLogListResp>(`/api/v1/performance-logs`, req)
}

/**
 * @description 
 * @param req
 */
export function performanceLogExport(req: components.PerformanceLogExportReq) {
	return webapi.get<components.PerformanceLogExportResp>(`/api/v1/performance-logs/export`, req)
}

/**
 * @description 
 */
export function sdkFileUpload() {
	return webapi.post<components.SdkFileUploadResp>(`/sdk/file/upload`)
}

/**
 * @description 
 * @param req
 */
export function sdkCallLogExport(req: components.SdkCallLogExportReq) {
	return webapi.get<null>(`/api/v1/sdk/call/log/export`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkCallLogList(req: components.SdkCallLogListReq) {
	return webapi.get<components.SdkCallLogListResp>(`/api/v1/sdk/call/log/list`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkInterfaceCreate(req: components.SdkInterfaceCreateReq) {
	return webapi.post<null>(`/api/v1/sdk/interface/create`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkInterfaceDelete(req: components.SdkInterfaceDeleteReq) {
	return webapi.post<null>(`/api/v1/sdk/interface/delete`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkInterfaceList(req: components.SdkInterfaceListReq) {
	return webapi.get<components.SdkInterfaceListResp>(`/api/v1/sdk/interface/list`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkInterfaceUpdate(req: components.SdkInterfaceUpdateReq) {
	return webapi.post<null>(`/api/v1/sdk/interface/update`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkApiKeyBindList(req: components.SdkApiKeyBindListReq) {
	return webapi.get<components.SdkApiKeyBindListResp>(`/api/v1/sdk/key/apis`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkApiKeyBindSave(req: components.SdkApiKeyBindSaveReq) {
	return webapi.post<null>(`/api/v1/sdk/key/apis/save`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkApiKeyCreate(req: components.SdkApiKeyCreateReq) {
	return webapi.post<components.SdkApiKeyCreateResp>(`/api/v1/sdk/key/create`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkApiKeyDelete(req: components.SdkApiKeyDeleteReq) {
	return webapi.post<null>(`/api/v1/sdk/key/delete`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkApiKeyList(req: components.SdkApiKeyListReq) {
	return webapi.get<components.SdkApiKeyListResp>(`/api/v1/sdk/key/list`, req)
}

/**
 * @description 
 * @param req
 */
export function sdkApiKeyUpdate(req: components.SdkApiKeyUpdateReq) {
	return webapi.post<null>(`/api/v1/sdk/key/update`, req)
}

/**
 * @description 
 * @param req
 */
export function configList(req: components.ConfigListReq) {
	return webapi.get<components.ConfigListResp>(`/api/v1/configs`, req)
}

/**
 * @description 
 * @param req
 */
export function configCreate(req: components.ConfigCreateReq) {
	return webapi.post<null>(`/api/v1/configs`, req)
}

/**
 * @description 
 * @param req
 */
export function configUpdate(req: components.ConfigUpdateReq) {
	return webapi.put<null>(`/api/v1/configs`, req)
}

/**
 * @description 
 * @param req
 */
export function configDelete(req: components.ConfigDeleteReq) {
	return webapi.delete<null>(`/api/v1/configs`, req)
}

/**
 * @description 
 * @param req
 */
export function configGet(req: components.ConfigGetReq) {
	return webapi.get<components.ConfigGetResp>(`/api/v1/configs/get`, req)
}

/**
 * @description 
 * @param req
 */
export function dictGet(req: components.DictGetReq) {
	return webapi.get<components.DictGetResp>(`/api/v1/dict`, req)
}

/**
 * @description 
 * @param req
 */
export function dictBatchGet(req: components.DictBatchGetReq) {
	return webapi.post<components.DictBatchGetResp>(`/api/v1/dict/batch`, req)
}

/**
 * @description 
 * @param req
 */
export function dictItemList(req: components.DictItemListReq) {
	return webapi.get<components.DictItemListResp>(`/api/v1/dict-items`, req)
}

/**
 * @description 
 * @param req
 */
export function dictItemCreate(req: components.DictItemCreateReq) {
	return webapi.post<null>(`/api/v1/dict-items`, req)
}

/**
 * @description 
 * @param req
 */
export function dictItemUpdate(req: components.DictItemUpdateReq) {
	return webapi.put<null>(`/api/v1/dict-items`, req)
}

/**
 * @description 
 * @param req
 */
export function dictItemDelete(req: components.DictItemDeleteReq) {
	return webapi.delete<null>(`/api/v1/dict-items`, req)
}

/**
 * @description 
 * @param req
 */
export function dictTypeList(req: components.DictTypeListReq) {
	return webapi.get<components.DictTypeListResp>(`/api/v1/dict-types`, req)
}

/**
 * @description 
 * @param req
 */
export function dictTypeCreate(req: components.DictTypeCreateReq) {
	return webapi.post<null>(`/api/v1/dict-types`, req)
}

/**
 * @description 
 * @param req
 */
export function dictTypeUpdate(req: components.DictTypeUpdateReq) {
	return webapi.put<null>(`/api/v1/dict-types`, req)
}

/**
 * @description 
 * @param req
 */
export function dictTypeDelete(req: components.DictTypeDeleteReq) {
	return webapi.delete<null>(`/api/v1/dict-types`, req)
}

/**
 * @description 
 * @param req
 */
export function fileList(req: components.FileListReq) {
	return webapi.get<components.FileListResp>(`/api/v1/files`, req)
}

/**
 * @description 
 * @param req
 */
export function fileCreate(req: components.FileCreateReq) {
	return webapi.post<null>(`/api/v1/files`, req)
}

/**
 * @description 
 * @param req
 */
export function fileUpdate(req: components.FileUpdateReq) {
	return webapi.put<null>(`/api/v1/files`, req)
}

/**
 * @description 
 * @param req
 */
export function fileDelete(req: components.FileDeleteReq) {
	return webapi.delete<null>(`/api/v1/files`, req)
}

/**
 * @description 
 * @param req
 */
export function fileDownload(req: components.FileDownloadReq) {
	return webapi.get<components.FileDownloadResp>(`/api/v1/files/download`, req)
}

/**
 * @description 
 */
export function fileUpload() {
	return webapi.post<components.FileUploadResp>(`/api/v1/files/upload`)
}

/**
 * @description 
 * @param req
 */
export function noticeList(req: components.NoticeListReq) {
	return webapi.get<components.NoticeListResp>(`/api/v1/notices`, req)
}

/**
 * @description 
 * @param req
 */
export function noticeCreate(req: components.NoticeCreateReq) {
	return webapi.post<components.Response>(`/api/v1/notices`, req)
}

/**
 * @description 
 * @param req
 */
export function noticeUpdate(req: components.NoticeUpdateReq) {
	return webapi.put<components.Response>(`/api/v1/notices`, req)
}

/**
 * @description 
 * @param req
 */
export function noticeDelete(req: components.NoticeDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/notices`, req)
}

/**
 * @description 
 * @param req
 */
export function notificationList(req: components.NotificationListReq) {
	return webapi.get<components.NotificationListResp>(`/api/v1/notifications`, req)
}

/**
 * @description 
 * @param req
 */
export function notificationDelete(req: components.NotificationDeleteReq) {
	return webapi.delete<components.Response>(`/api/v1/notifications`, req)
}

/**
 * @description 
 * @param req
 */
export function notificationRead(req: components.NotificationReadReq) {
	return webapi.put<components.Response>(`/api/v1/notifications/read`, req)
}

/**
 * @description 
 */
export function notificationClearRead() {
	return webapi.delete<components.Response>(`/api/v1/notifications/read`)
}

/**
 * @description 
 */
export function notificationReadAll() {
	return webapi.put<components.Response>(`/api/v1/notifications/read-all`)
}

/**
 * @description 
 * @param req
 */
export function taskRecent(req: components.TaskRecentReq) {
	return webapi.get<components.TaskRecentResp>(`/api/v1/tasks/recent`, req)
}

/**
 * @description 
 * @param req
 */
export function taskList(req: components.TaskListReq) {
	return webapi.get<components.TaskListResp>(`/api/v1/tasks`, req)
}

/**
 * @description 
 * @param req
 */
export function taskCancel(req: components.TaskCancelReq) {
	return webapi.post<components.Response>(`/api/v1/tasks/cancel`, req)
}

/**
 * @description 
 * @param req
 */
export function taskDetail(req: components.TaskDetailReq) {
	return webapi.get<components.TaskDetailResp>(`/api/v1/tasks/detail`, req)
}

/**
 * @description 
 * @param req
 */
export function m3u8Proxy(req: components.M3u8ProxyReq) {
	return webapi.get<null>(`/api/v1/m3u8/proxy`, req)
}

/**
 * @description 
 * @param req
 */
export function publicVideoDetail(req: components.PublicVideoDetailReq) {
	return webapi.get<components.PublicVideoDetailResp>(`/api/v1/public/videos/info`, req)
}

/**
 * @description 
 * @param req
 */
export function publicVideoList(req: components.PublicVideoListReq) {
	return webapi.get<components.PublicVideoListResp>(`/api/v1/public/videos/list`, req)
}

/**
 * @description 
 * @param req
 */
export function videoList(req: components.VideoListReq) {
	return webapi.get<components.VideoListResp>(`/api/v1/videos`, req)
}

/**
 * @description 
 * @param req
 */
export function videoCreate(req: components.VideoCreateReq) {
	return webapi.post<null>(`/api/v1/videos`, req)
}

/**
 * @description 
 * @param req
 */
export function videoUpdate(req: components.VideoUpdateReq) {
	return webapi.put<null>(`/api/v1/videos`, req)
}

/**
 * @description 
 * @param req
 */
export function videoDelete(req: components.VideoDeleteReq) {
	return webapi.delete<null>(`/api/v1/videos`, req)
}

/**
 * @description 
 * @param req
 */
export function videoCollect(req: components.VideoCollectReq) {
	return webapi.post<components.VideoCollectResp>(`/api/v1/videos/collect`, req)
}

/**
 * @description 
 */
export function videoCollectOptions() {
	return webapi.options<null>(`/api/v1/videos/collect`)
}
