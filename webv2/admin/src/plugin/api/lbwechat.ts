import type * as lbwechat from "./model/lbwechat";
import api from "./config/api";

async function handleWxGzhAuth(req: lbwechat.HandleWxGzhAuthReq): Promise<lbwechat.HandleWxGzhAuthRsp> {
    const response = await api.get<lbwechat.HandleWxGzhAuthRsp>("/lbwechat/HandleWxGzhAuth", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function handleWxGzhMsg(req: lbwechat.HandleWxGzhMsgReq): Promise<lbwechat.HandleWxGzhMsgRsp> {
    const response = await api.post<lbwechat.HandleWxGzhMsgRsp>("/lbwechat/HandleWxGzhMsg", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

export default {
	handleWxGzhAuth,
	handleWxGzhMsg
}
