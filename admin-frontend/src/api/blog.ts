import {
  blogTagList,
  blogTagOptions,
  blogTagCreate,
  blogTagUpdate,
  blogTagDelete,
  blogArticleList,
  blogArticleCreate,
  blogArticleUpdate,
  blogArticleDelete,
  blogArticleDetail,
  blogArticleSubmit,
  blogArticlePublish,
  blogArticleUnpublish,
  blogArticleAudit,
  blogArticleAuditUnpublish,
  publicBlogArticleList,
  publicBlogArticleDetail
} from '@/api/generated/admin';
import type {
  BlogTagListReq,
  BlogTagListResp,
  BlogTagOptionsReq,
  BlogTagOptionsResp,
  BlogTagCreateReq,
  BlogTagUpdateReq,
  BlogTagDeleteReq,
  BlogArticleListReq,
  BlogArticleListResp,
  BlogArticleCreateReq,
  BlogArticleUpdateReq,
  BlogArticleDeleteReq,
  BlogArticleDetailReq,
  BlogArticleDetailResp,
  BlogArticleSubmitReq,
  BlogArticlePublishReq,
  BlogArticleUnpublishReq,
  BlogArticleAuditReq,
  BlogArticleAuditUnpublishReq,
  PublicBlogArticleListReq,
  PublicBlogArticleListResp,
  PublicBlogArticleDetailReq,
  PublicBlogArticleDetailResp,
  Response
} from '@/api/generated/admin';

// 标签管理
export const blogApi = {
  tagList: (req: BlogTagListReq) => blogTagList(req) as Promise<BlogTagListResp>,
  tagOptions: (req?: BlogTagOptionsReq) => blogTagOptions(req || {}) as Promise<BlogTagOptionsResp>,
  tagCreate: (req: BlogTagCreateReq) => blogTagCreate(req),
  tagUpdate: (req: BlogTagUpdateReq) => blogTagUpdate(req),
  tagDelete: (req: BlogTagDeleteReq) => blogTagDelete(req),

  // 文章管理（后台）
  articleList: (req: BlogArticleListReq) => blogArticleList(req) as Promise<BlogArticleListResp>,
  articleCreate: (req: BlogArticleCreateReq) => blogArticleCreate(req) as Promise<Response>,
  articleUpdate: (req: BlogArticleUpdateReq) => blogArticleUpdate(req) as Promise<Response>,
  articleDelete: (req: BlogArticleDeleteReq) => blogArticleDelete(req) as Promise<Response>,
  articleDetail: (req: BlogArticleDetailReq) => blogArticleDetail(req) as Promise<BlogArticleDetailResp>,
  articleSubmit: (req: BlogArticleSubmitReq) => blogArticleSubmit(req) as Promise<Response>,
  articlePublish: (req: BlogArticlePublishReq) => blogArticlePublish(req) as Promise<Response>,
  articleUnpublish: (req: BlogArticleUnpublishReq) => blogArticleUnpublish(req) as Promise<Response>,

  // 审核操作
  articleAudit: (req: BlogArticleAuditReq) => blogArticleAudit(req) as Promise<Response>,
  articleAuditUnpublish: (req: BlogArticleAuditUnpublishReq) => blogArticleAuditUnpublish(req) as Promise<Response>,

  // 公共文章接口
  publicList: (req: PublicBlogArticleListReq) =>
    publicBlogArticleList(req) as Promise<PublicBlogArticleListResp>,
  publicDetail: (req: PublicBlogArticleDetailReq) =>
    publicBlogArticleDetail(req) as Promise<PublicBlogArticleDetailResp>
};

