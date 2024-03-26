import api from "./config/api";

export default class lbbillApi {
    async addBillSys(req: lbbill.AddBillSysReq): Promise<lbbill.AddBillSysRsp> {
        const response = await api.post<lbbill.AddBillSysRsp>("/lbbill/AddBillSys", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async delBillSysList(req: lbbill.DelBillSysListReq): Promise<lbbill.DelBillSysListRsp> {
        const response = await api.post<lbbill.DelBillSysListRsp>("/lbbill/DelBillSysList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateBillSys(req: lbbill.UpdateBillSysReq): Promise<lbbill.UpdateBillSysRsp> {
        const response = await api.post<lbbill.UpdateBillSysRsp>("/lbbill/UpdateBillSys", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getBillSys(req: lbbill.GetBillSysReq): Promise<lbbill.GetBillSysRsp> {
        const response = await api.post<lbbill.GetBillSysRsp>("/lbbill/GetBillSys", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getBillSysList(req: lbbill.GetBillSysListReq): Promise<lbbill.GetBillSysListRsp> {
        const response = await api.post<lbbill.GetBillSysListRsp>("/lbbill/GetBillSysList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async addBillCategorySys(req: lbbill.AddBillCategorySysReq): Promise<lbbill.AddBillCategorySysRsp> {
        const response = await api.post<lbbill.AddBillCategorySysRsp>("/lbbill/AddBillCategorySys", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async delBillCategorySysList(req: lbbill.DelBillCategorySysListReq): Promise<lbbill.DelBillCategorySysListRsp> {
        const response = await api.post<lbbill.DelBillCategorySysListRsp>("/lbbill/DelBillCategorySysList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateBillCategorySys(req: lbbill.UpdateBillCategorySysReq): Promise<lbbill.UpdateBillCategorySysRsp> {
        const response = await api.post<lbbill.UpdateBillCategorySysRsp>("/lbbill/UpdateBillCategorySys", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getBillCategorySys(req: lbbill.GetBillCategorySysReq): Promise<lbbill.GetBillCategorySysRsp> {
        const response = await api.post<lbbill.GetBillCategorySysRsp>("/lbbill/GetBillCategorySys", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getBillCategorySysList(req: lbbill.GetBillCategorySysListReq): Promise<lbbill.GetBillCategorySysListRsp> {
        const response = await api.post<lbbill.GetBillCategorySysListRsp>("/lbbill/GetBillCategorySysList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

}
const lbbillSingle = new lbbillApi();

export {
    lbbillSingle
}
