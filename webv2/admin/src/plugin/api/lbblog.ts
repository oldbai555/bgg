import api from "./config/api";

export default class lbblogApi {
    async addArticle(req: lbblog.AddArticleReq): Promise<lbblog.AddArticleRsp> {
        const response = await api.post<lbblog.AddArticleRsp>("/lbblog/AddArticle", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async delArticleList(req: lbblog.DelArticleListReq): Promise<lbblog.DelArticleListRsp> {
        const response = await api.post<lbblog.DelArticleListRsp>("/lbblog/DelArticleList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateArticle(req: lbblog.UpdateArticleReq): Promise<lbblog.UpdateArticleRsp> {
        const response = await api.post<lbblog.UpdateArticleRsp>("/lbblog/UpdateArticle", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getArticle(req: lbblog.GetArticleReq): Promise<lbblog.GetArticleRsp> {
        const response = await api.post<lbblog.GetArticleRsp>("/lbblog/GetArticle", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getArticleList(req: lbblog.GetArticleListReq): Promise<lbblog.GetArticleListRsp> {
        const response = await api.post<lbblog.GetArticleListRsp>("/lbblog/GetArticleList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async addCategory(req: lbblog.AddCategoryReq): Promise<lbblog.AddCategoryRsp> {
        const response = await api.post<lbblog.AddCategoryRsp>("/lbblog/AddCategory", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async delCategoryList(req: lbblog.DelCategoryListReq): Promise<lbblog.DelCategoryListRsp> {
        const response = await api.post<lbblog.DelCategoryListRsp>("/lbblog/DelCategoryList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateCategory(req: lbblog.UpdateCategoryReq): Promise<lbblog.UpdateCategoryRsp> {
        const response = await api.post<lbblog.UpdateCategoryRsp>("/lbblog/UpdateCategory", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getCategory(req: lbblog.GetCategoryReq): Promise<lbblog.GetCategoryRsp> {
        const response = await api.post<lbblog.GetCategoryRsp>("/lbblog/GetCategory", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getCategoryList(req: lbblog.GetCategoryListReq): Promise<lbblog.GetCategoryListRsp> {
        const response = await api.post<lbblog.GetCategoryListRsp>("/lbblog/GetCategoryList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async addComment(req: lbblog.AddCommentReq): Promise<lbblog.AddCommentRsp> {
        const response = await api.post<lbblog.AddCommentRsp>("/lbblog/AddComment", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async delCommentList(req: lbblog.DelCommentListReq): Promise<lbblog.DelCommentListRsp> {
        const response = await api.post<lbblog.DelCommentListRsp>("/lbblog/DelCommentList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async updateComment(req: lbblog.UpdateCommentReq): Promise<lbblog.UpdateCommentRsp> {
        const response = await api.post<lbblog.UpdateCommentRsp>("/lbblog/UpdateComment", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getComment(req: lbblog.GetCommentReq): Promise<lbblog.GetCommentRsp> {
        const response = await api.post<lbblog.GetCommentRsp>("/lbblog/GetComment", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

    async getCommentList(req: lbblog.GetCommentListReq): Promise<lbblog.GetCommentListRsp> {
        const response = await api.post<lbblog.GetCommentListRsp>("/lbblog/GetCommentList", req);
        if (response.code !== 200) {
            return Promise.reject(response.message);
        }
        return response.data;
    }

}
const lbblogSingle = new lbblogApi();

export {
    lbblogSingle
}