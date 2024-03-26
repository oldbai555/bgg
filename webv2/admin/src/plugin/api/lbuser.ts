import api from "./config/api";

export default class lbuserApi {
    async login(req: lbuser.LoginReq): Promise<lbuser.LoginRsp> {
        const response = await api.post<lbuser.LoginRsp>("/lbuser/Login", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async logout(req: lbuser.LogoutReq): Promise<lbuser.LogoutRsp> {
        const response = await api.post<lbuser.LogoutRsp>("/lbuser/Logout", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getLoginUser(req: lbuser.GetLoginUserReq): Promise<lbuser.GetLoginUserRsp> {
        const response = await api.post<lbuser.GetLoginUserRsp>("/lbuser/GetLoginUser", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateLoginUserInfo(req: lbuser.UpdateLoginUserInfoReq): Promise<lbuser.UpdateLoginUserInfoRsp> {
        const response = await api.post<lbuser.UpdateLoginUserInfoRsp>("/lbuser/UpdateLoginUserInfo", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async resetPassword(req: lbuser.ResetPasswordReq): Promise<lbuser.ResetPasswordRsp> {
        const response = await api.post<lbuser.ResetPasswordRsp>("/lbuser/ResetPassword", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getFrontUser(req: lbuser.GetFrontUserReq): Promise<lbuser.GetLoginUserRsp> {
        const response = await api.post<lbuser.GetLoginUserRsp>("/lbuser/GetFrontUser", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async addUser(req: lbuser.AddUserReq): Promise<lbuser.AddUserRsp> {
        const response = await api.post<lbuser.AddUserRsp>("/lbuser/AddUser", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async delUserList(req: lbuser.DelUserListReq): Promise<lbuser.DelUserListRsp> {
        const response = await api.post<lbuser.DelUserListRsp>("/lbuser/DelUserList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateUser(req: lbuser.UpdateUserReq): Promise<lbuser.UpdateUserRsp> {
        const response = await api.post<lbuser.UpdateUserRsp>("/lbuser/UpdateUser", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getUser(req: lbuser.GetUserReq): Promise<lbuser.GetUserRsp> {
        const response = await api.post<lbuser.GetUserRsp>("/lbuser/GetUser", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getUserList(req: lbuser.GetUserListReq): Promise<lbuser.GetUserListRsp> {
        const response = await api.post<lbuser.GetUserListRsp>("/lbuser/GetUserList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }
}
const lbuserSingle = new lbuserApi();

export {
    lbuserSingle
}