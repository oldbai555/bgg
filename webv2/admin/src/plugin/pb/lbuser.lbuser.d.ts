// Code generated by rpc_gen. DO NOT EDIT.

declare namespace lbuser {
    export const enum ErrCode {
        Nil = 0,
        ErrUserNotFound = 30001, //用户找不到
        ErrPasswordInvalid = 30002, //密码无效
        ErrUserExit = 30003, // 用户重复存在
        ErrUserLoginExpired = 30004, // 登陆信息过期
        ErrGrpcParseContextFail = 30005, // 无效请求
        ErrLoginUserNotEqualChangeUser = 30006, // 当前用户无法操作指定用户
        ErrUserNameInvalid = 30007, //用户名无效
    }

    export const enum ModelUser_Role {
        RoleNil = 0,
        RoleAdmin = 1,
    }

    export interface ModelUser {
        id?: string;
        created_at?: number;
        updated_at?: number;
        deleted_at?: number;
        username?: string;
        password?: string;
        avatar?: string;
        nickname?: string;
        email?: string;
        github?: string;
        desc?: string;
        role?: number;
    }

    export interface BaseUser {
        id?: string;
        username?: string;
        avatar?: string;
        nickname?: string;
        email?: string;
        github?: string;
        desc?: string;
        role?: number;
    }

    export interface LoginReq {
        username?: string;
        password?: string;
    }

    export interface LoginRsp {
        sid?: string;
    }

    export interface LogoutReq {
        sid?: string;
    }

    export interface LogoutRsp {
    }

    export interface GetLoginUserReq {
        sid?: string;
    }

    export interface GetLoginUserRsp {
        base_user?: BaseUser;
    }

    export interface UpdateLoginUserInfoReq {
        user?: BaseUser;
    }

    export interface UpdateLoginUserInfoRsp {
    }

    export interface ResetPasswordReq {
        old_password?: string;
        new_password?: string;
    }

    export interface ResetPasswordRsp {
    }

    export interface GetFrontUserReq {
    }

    export interface AddUserReq {
        data?: ModelUser;
    }

    export interface AddUserRsp {
        data?: ModelUser;
    }

    export interface UpdateUserReq {
        data?: ModelUser;
    }

    export interface UpdateUserRsp {
    }

    export interface DelUserListReq {
        // @ref_to: GetUserListReq.ListOption
        list_option?: lb.ListOption;
    }

    export interface DelUserListRsp {
    }

    export interface GetUserReq {
        id?: string;
    }

    export interface GetUserRsp {
        data?: ModelUser;
    }

    export const enum GetUserListReq_ListOption {
        ListOptionNil = 0,
        ListOptionLikeUsername = 1,
    }

    export interface GetUserListReq {
        list_option?: lb.ListOption;
    }

    export interface GetUserListRsp {
        paginate?: lb.Paginate;
        list?: Array<ModelUser>;
    }

    export interface lbuserService {
        // @desc: 登录
        Login<R extends LoginReq, T>(r: R, o?: T): Promise<LoginRsp>;

        // @desc: 登出
        Logout<R extends LogoutReq, T>(r: R, o?: T): Promise<LogoutRsp>;

        // @desc: 获取登录用户的信息
        GetLoginUser<R extends GetLoginUserReq, T>(r: R, o?: T): Promise<GetLoginUserRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error: 更新登陆的用户信息
        UpdateLoginUserInfo<R extends UpdateLoginUserInfoReq, T>(r: R, o?: T): Promise<UpdateLoginUserInfoRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error:
        ResetPassword<R extends ResetPasswordReq, T>(r: R, o?: T): Promise<ResetPasswordRsp>;

        // @cat: front
        // @name:
        // @desc:
        // @error:
        GetFrontUser<R extends GetFrontUserReq, T>(r: R, o?: T): Promise<GetLoginUserRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error:
        AddUser<R extends AddUserReq, T>(r: R, o?: T): Promise<AddUserRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error:
        DelUserList<R extends DelUserListReq, T>(r: R, o?: T): Promise<DelUserListRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error:
        UpdateUser<R extends UpdateUserReq, T>(r: R, o?: T): Promise<UpdateUserRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error:
        GetUser<R extends GetUserReq, T>(r: R, o?: T): Promise<GetUserRsp>;

        // @cat:
        // @name:
        // @desc:
        // @error:
        GetUserList<R extends GetUserListReq, T>(r: R, o?: T): Promise<GetUserListRsp>;
    }
}