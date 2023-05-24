import type * as lbblog from "./model/lbblog";
import api from "./config/api";

async function getArticleList(req: lbblog.GetArticleListReq): Promise<lbblog.GetArticleListRsp> {
    const response = await api.post<lbblog.GetArticleListRsp>("/lbblog/GetArticleList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getArticle(req: lbblog.GetArticleReq): Promise<lbblog.GetArticleRsp> {
    const response = await api.post<lbblog.GetArticleRsp>("/lbblog/GetArticle", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateArticle(req: lbblog.UpdateArticleReq): Promise<lbblog.UpdateArticleRsp> {
    const response = await api.post<lbblog.UpdateArticleRsp>("/lbblog/UpdateArticle", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function delArticle(req: lbblog.DelArticleReq): Promise<lbblog.DelArticleRsp> {
    const response = await api.post<lbblog.DelArticleRsp>("/lbblog/DelArticle", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function addArticle(req: lbblog.AddArticleReq): Promise<lbblog.AddArticleRsp> {
    const response = await api.post<lbblog.AddArticleRsp>("/lbblog/AddArticle", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getCategoryList(req: lbblog.GetCategoryListReq): Promise<lbblog.GetCategoryListRsp> {
    const response = await api.post<lbblog.GetCategoryListRsp>("/lbblog/GetCategoryList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getCategory(req: lbblog.GetCategoryReq): Promise<lbblog.GetCategoryRsp> {
    const response = await api.post<lbblog.GetCategoryRsp>("/lbblog/GetCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateCategory(req: lbblog.UpdateCategoryReq): Promise<lbblog.UpdateCategoryRsp> {
    const response = await api.post<lbblog.UpdateCategoryRsp>("/lbblog/UpdateCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function delCategory(req: lbblog.DelCategoryReq): Promise<lbblog.DelCategoryRsp> {
    const response = await api.post<lbblog.DelCategoryRsp>("/lbblog/DelCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function addCategory(req: lbblog.AddCategoryReq): Promise<lbblog.AddCategoryRsp> {
    const response = await api.post<lbblog.AddCategoryRsp>("/lbblog/AddCategory", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getCommentList(req: lbblog.GetCommentListReq): Promise<lbblog.GetCommentListRsp> {
    const response = await api.post<lbblog.GetCommentListRsp>("/lbblog/GetCommentList", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function getComment(req: lbblog.GetCommentReq): Promise<lbblog.GetCommentRsp> {
    const response = await api.post<lbblog.GetCommentRsp>("/lbblog/GetComment", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function updateComment(req: lbblog.UpdateCommentReq): Promise<lbblog.UpdateCommentRsp> {
    const response = await api.post<lbblog.UpdateCommentRsp>("/lbblog/UpdateComment", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function delComment(req: lbblog.DelCommentReq): Promise<lbblog.DelCommentRsp> {
    const response = await api.post<lbblog.DelCommentRsp>("/lbblog/DelComment", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

async function addComment(req: lbblog.AddCommentReq): Promise<lbblog.AddCommentRsp> {
    const response = await api.post<lbblog.AddCommentRsp>("/lbblog/AddComment", req);
    if (response.code !== 200) {
        return Promise.reject(response.message);
    }
    return response.data;
}

export default {
	getArticleList,
	getArticle,
	updateArticle,
	delArticle,
	addArticle,
	getCategoryList,
	getCategory,
	updateCategory,
	delCategory,
	addCategory,
	getCommentList,
	getComment,
	updateComment,
	delComment,
	addComment
}
