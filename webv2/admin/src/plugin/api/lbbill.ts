import type * as lbbill from "./model/lbbill";
import api from "./config/api";

async function addBill(req: lbbill.AddBillReq): Promise<lbbill.AddBillRsp> {
    const response = await api.post<lbbill.AddBillRsp>("/lbbill/AddBill", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function delBill(req: lbbill.DelBillReq): Promise<lbbill.DelBillRsp> {
    const response = await api.post<lbbill.DelBillRsp>("/lbbill/DelBill", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateBill(req: lbbill.UpdateBillReq): Promise<lbbill.UpdateBillRsp> {
    const response = await api.post<lbbill.UpdateBillRsp>("/lbbill/UpdateBill", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getBill(req: lbbill.GetBillReq): Promise<lbbill.GetBillRsp> {
    const response = await api.post<lbbill.GetBillRsp>("/lbbill/GetBill", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getBillList(req: lbbill.GetBillListReq): Promise<lbbill.GetBillListRsp> {
    const response = await api.post<lbbill.GetBillListRsp>("/lbbill/GetBillList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function addBillCategory(req: lbbill.AddBillCategoryReq): Promise<lbbill.AddBillCategoryRsp> {
    const response = await api.post<lbbill.AddBillCategoryRsp>("/lbbill/AddBillCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function delBillCategory(req: lbbill.DelBillCategoryReq): Promise<lbbill.DelBillCategoryRsp> {
    const response = await api.post<lbbill.DelBillCategoryRsp>("/lbbill/DelBillCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateBillCategory(req: lbbill.UpdateBillCategoryReq): Promise<lbbill.UpdateBillCategoryRsp> {
    const response = await api.post<lbbill.UpdateBillCategoryRsp>("/lbbill/UpdateBillCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getBillCategory(req: lbbill.GetBillCategoryReq): Promise<lbbill.GetBillCategoryRsp> {
    const response = await api.post<lbbill.GetBillCategoryRsp>("/lbbill/GetBillCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getBillCategoryList(req: lbbill.GetBillCategoryListReq): Promise<lbbill.GetBillCategoryListRsp> {
    const response = await api.post<lbbill.GetBillCategoryListRsp>("/lbbill/GetBillCategoryList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

export default {
	addBill,
	delBill,
	updateBill,
	getBill,
	getBillList,
	addBillCategory,
	delBillCategory,
	updateBillCategory,
	getBillCategory,
	getBillCategoryList
}
