import type * as lbuser from "./model/lbuser";
import api from "./config/api";

async function login(req: lbuser.LoginReq): Promise<lbuser.LoginRsp> {
    const response = await api.post<lbuser.LoginRsp>("/lbuser/Login", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function logout(req: lbuser.LogoutReq): Promise<lbuser.LogoutRsp> {
    const response = await api.post<lbuser.LogoutRsp>("/lbuser/Logout", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getLoginUser(req: lbuser.GetLoginUserReq): Promise<lbuser.GetLoginUserRsp> {
    const response = await api.get<lbuser.GetLoginUserRsp>("/lbuser/GetLoginUser", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateLoginUserInfo(req: lbuser.UpdateLoginUserInfoReq): Promise<lbuser.UpdateLoginUserInfoRsp> {
    const response = await api.post<lbuser.UpdateLoginUserInfoRsp>("/lbuser/UpdateLoginUserInfo", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function addUser(req: lbuser.AddUserReq): Promise<lbuser.AddUserRsp> {
    const response = await api.post<lbuser.AddUserRsp>("/lbuser/AddUser", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getUserList(req: lbuser.GetUserListReq): Promise<lbuser.GetUserListRsp> {
    const response = await api.post<lbuser.GetUserListRsp>("/lbuser/GetUserList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function delUser(req: lbuser.DelUserReq): Promise<lbuser.DelUserRsp> {
    const response = await api.post<lbuser.DelUserRsp>("/lbuser/DelUser", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getUser(req: lbuser.GetUserReq): Promise<lbuser.GetUserRsp> {
    const response = await api.post<lbuser.GetUserRsp>("/lbuser/GetUser", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateUserNameWithRole(req: lbuser.UpdateUserNameWithRoleReq): Promise<lbuser.UpdateUserNameWithRoleRsp> {
    const response = await api.post<lbuser.UpdateUserNameWithRoleRsp>("/lbuser/UpdateUserNameWithRole", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function resetPassword(req: lbuser.ResetPasswordReq): Promise<lbuser.ResetPasswordRsp> {
    const response = await api.post<lbuser.ResetPasswordRsp>("/lbuser/ResetPassword", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getFrontUser(req: lbuser.GetFrontUserReq): Promise<lbuser.GetLoginUserRsp> {
    const response = await api.get<lbuser.GetLoginUserRsp>("/lbuser/GetFrontUser", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

export default {
	login,
	logout,
	getLoginUser,
	updateLoginUserInfo,
	addUser,
	getUserList,
	delUser,
	getUser,
	updateUserNameWithRole,
	resetPassword,
	getFrontUser
}
