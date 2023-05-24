import type * as lbchatgpt from "./model/lbchatgpt";
import api from "./config/api";

async function chatCompletion(req: lbchatgpt.ChatCompletionReq): Promise<lbchatgpt.ChatCompletionRsp> {
    const response = await api.post<lbchatgpt.ChatCompletionRsp>("/lbchatgpt/ChatCompletion", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

export default {
	chatCompletion
}
