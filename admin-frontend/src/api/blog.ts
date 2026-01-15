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
  blogArticleTop,
  blogArticleUntop,
  publicBlogArticleList,
  publicBlogArticleDetail,
  publicBlogFriendLinkList,
  publicBlogSocialInfoList,
  publicBlogTagList,
  publicBlogAuthorInfo,
  publicBlogArticleStats,
  publicBlogArticlePrev,
  publicBlogArticleNext
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
  BlogArticleTopReq,
  BlogArticleUntopReq,
  PublicBlogArticleListReq,
  PublicBlogArticleListResp,
  PublicBlogArticleDetailReq,
  PublicBlogArticleDetailResp,
  PublicBlogFriendLinkListResp,
  PublicBlogSocialInfoListResp,
  PublicBlogTagListResp,
  PublicBlogAuthorInfoResp,
  PublicBlogArticleStatsResp,
  PublicBlogArticlePrevReq,
  PublicBlogArticlePrevResp,
  PublicBlogArticleNextReq,
  PublicBlogArticleNextResp,
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

  // 置顶操作
  articleTop: (req: BlogArticleTopReq) => blogArticleTop(req) as Promise<Response>,
  articleUntop: (req: BlogArticleUntopReq) => blogArticleUntop(req) as Promise<Response>,

  // 公共文章接口
  publicList: (req: PublicBlogArticleListReq) =>
    publicBlogArticleList(req) as Promise<PublicBlogArticleListResp>,
  publicDetail: (req: PublicBlogArticleDetailReq) =>
    publicBlogArticleDetail(req) as Promise<PublicBlogArticleDetailResp>,

  // 公共友情链接接口
  publicFriendLinkList: () =>
    publicBlogFriendLinkList() as Promise<PublicBlogFriendLinkListResp>,

  // 公共社交信息接口
  publicSocialInfoList: () =>
    publicBlogSocialInfoList() as Promise<PublicBlogSocialInfoListResp>,

  // 公共标签列表接口
  publicTagList: () =>
    publicBlogTagList() as Promise<PublicBlogTagListResp>,

  // 公共作者信息接口
  publicAuthorInfo: () =>
    publicBlogAuthorInfo() as Promise<PublicBlogAuthorInfoResp>,

  // 公共文章统计接口
  publicArticleStats: () =>
    publicBlogArticleStats() as Promise<PublicBlogArticleStatsResp>,

  // 公共相邻文章接口
  publicArticlePrev: (req: PublicBlogArticlePrevReq) =>
    publicBlogArticlePrev(req) as Promise<PublicBlogArticlePrevResp>,
  publicArticleNext: (req: PublicBlogArticleNextReq) =>
    publicBlogArticleNext(req) as Promise<PublicBlogArticleNextResp>
};

