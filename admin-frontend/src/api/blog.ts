import {
  // 公开接口
  publicBlogArticleList,
  publicBlogArticleDetail,
  publicBlogFriendLinkList,
  publicBlogSocialInfoList,
  publicBlogTagList,
  publicBlogAuthorInfo,
  publicBlogArticleStats,
  publicBlogArticlePrev,
  publicBlogArticleNext,
  // 后台管理接口
  blogArticleList,
  blogArticleCreate,
  blogArticleUpdate,
  blogArticleDelete,
  blogArticleDetail,
  blogArticlePublish,
  blogArticleSubmit,
  blogArticleTop,
  blogArticleUnpublish,
  blogArticleUntop,
  blogArticleAudit,
  blogArticleAuditUnpublish,
  blogTagList,
  blogTagCreate,
  blogTagUpdate,
  blogTagDelete,
  blogTagOptions,
  blogFriendLinkList,
  blogFriendLinkCreate,
  blogFriendLinkUpdate,
  blogFriendLinkDelete,
  blogSocialInfoList,
  blogSocialInfoCreate,
  blogSocialInfoUpdate,
  blogSocialInfoDelete
} from '@/api/generated/admin'
import type {
  // 公开接口类型
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
  // 后台管理接口类型
  BlogArticleListReq,
  BlogArticleListResp,
  BlogArticleCreateReq,
  BlogArticleUpdateReq,
  BlogArticleDeleteReq,
  BlogArticleDetailReq,
  BlogArticleDetailResp,
  BlogArticlePublishReq,
  BlogArticleSubmitReq,
  BlogArticleTopReq,
  BlogArticleUnpublishReq,
  BlogArticleUntopReq,
  BlogArticleAuditReq,
  BlogArticleAuditUnpublishReq,
  BlogTagListReq,
  BlogTagListResp,
  BlogTagCreateReq,
  BlogTagUpdateReq,
  BlogTagDeleteReq,
  BlogTagOptionsReq,
  BlogTagOptionsResp,
  BlogFriendLinkListReq,
  BlogFriendLinkListResp,
  BlogFriendLinkCreateReq,
  BlogFriendLinkUpdateReq,
  BlogFriendLinkDeleteReq,
  BlogSocialInfoListReq,
  BlogSocialInfoListResp,
  BlogSocialInfoCreateReq,
  BlogSocialInfoUpdateReq,
  BlogSocialInfoDeleteReq
} from '@/api/generated/admin'

/**
 * 博客相关 API 封装
 */
export const blogApi = {
  // ========== 公开接口 ==========
  // 公开博客文章列表
  publicArticleList: (req: PublicBlogArticleListReq) =>
    publicBlogArticleList(req) as Promise<PublicBlogArticleListResp>,

  // 公开博客文章详情
  publicArticleDetail: (req: PublicBlogArticleDetailReq) =>
    publicBlogArticleDetail(req) as Promise<PublicBlogArticleDetailResp>,

  // 公开友情链接列表
  publicFriendLinkList: () =>
    publicBlogFriendLinkList() as Promise<PublicBlogFriendLinkListResp>,

  // 公开社交信息列表
  publicSocialInfoList: () =>
    publicBlogSocialInfoList() as Promise<PublicBlogSocialInfoListResp>,

  // 公开标签列表
  publicTagList: () =>
    publicBlogTagList() as Promise<PublicBlogTagListResp>,

  // 公开作者信息
  publicAuthorInfo: () =>
    publicBlogAuthorInfo() as Promise<PublicBlogAuthorInfoResp>,

  // 公开文章统计
  publicArticleStats: () =>
    publicBlogArticleStats() as Promise<PublicBlogArticleStatsResp>,

  // 上一篇文章
  publicArticlePrev: (req: PublicBlogArticlePrevReq) =>
    publicBlogArticlePrev(req) as Promise<PublicBlogArticlePrevResp>,

  // 下一篇文章
  publicArticleNext: (req: PublicBlogArticleNextReq) =>
    publicBlogArticleNext(req) as Promise<PublicBlogArticleNextResp>,

  // ========== 后台管理接口 ==========
  // 文章列表
  articleList: (req: BlogArticleListReq) =>
    blogArticleList(req) as Promise<BlogArticleListResp>,

  // 文章创建
  articleCreate: (req: BlogArticleCreateReq) =>
    blogArticleCreate(req),

  // 文章更新
  articleUpdate: (req: BlogArticleUpdateReq) =>
    blogArticleUpdate(req),

  // 文章删除
  articleDelete: (req: BlogArticleDeleteReq) =>
    blogArticleDelete(req),

  // 文章详情
  articleDetail: (req: BlogArticleDetailReq) =>
    blogArticleDetail(req) as Promise<BlogArticleDetailResp>,

  // 文章发布
  articlePublish: (req: BlogArticlePublishReq) =>
    blogArticlePublish(req),

  // 文章提交审核
  articleSubmit: (req: BlogArticleSubmitReq) =>
    blogArticleSubmit(req),

  // 文章置顶
  articleTop: (req: BlogArticleTopReq) =>
    blogArticleTop(req),

  // 文章取消发布
  articleUnpublish: (req: BlogArticleUnpublishReq) =>
    blogArticleUnpublish(req),

  // 文章取消置顶
  articleUntop: (req: BlogArticleUntopReq) =>
    blogArticleUntop(req),

  // 文章审核
  articleAudit: (req: BlogArticleAuditReq) =>
    blogArticleAudit(req),

  // 文章审核取消发布
  articleAuditUnpublish: (req: BlogArticleAuditUnpublishReq) =>
    blogArticleAuditUnpublish(req),

  // 标签列表
  tagList: (req: BlogTagListReq) =>
    blogTagList(req) as Promise<BlogTagListResp>,

  // 标签创建
  tagCreate: (req: BlogTagCreateReq) =>
    blogTagCreate(req),

  // 标签更新
  tagUpdate: (req: BlogTagUpdateReq) =>
    blogTagUpdate(req),

  // 标签删除
  tagDelete: (req: BlogTagDeleteReq) =>
    blogTagDelete(req),

  // 标签选项（用于下拉选择）
  tagOptions: (req: BlogTagOptionsReq) =>
    blogTagOptions(req) as Promise<BlogTagOptionsResp>,

  // 友情链接列表
  friendLinkList: (req: BlogFriendLinkListReq) =>
    blogFriendLinkList(req) as Promise<BlogFriendLinkListResp>,

  // 友情链接创建
  friendLinkCreate: (req: BlogFriendLinkCreateReq) =>
    blogFriendLinkCreate(req),

  // 友情链接更新
  friendLinkUpdate: (req: BlogFriendLinkUpdateReq) =>
    blogFriendLinkUpdate(req),

  // 友情链接删除
  friendLinkDelete: (req: BlogFriendLinkDeleteReq) =>
    blogFriendLinkDelete(req),

  // 社交信息列表
  socialInfoList: (req: BlogSocialInfoListReq) =>
    blogSocialInfoList(req) as Promise<BlogSocialInfoListResp>,

  // 社交信息创建
  socialInfoCreate: (req: BlogSocialInfoCreateReq) =>
    blogSocialInfoCreate(req),

  // 社交信息更新
  socialInfoUpdate: (req: BlogSocialInfoUpdateReq) =>
    blogSocialInfoUpdate(req),

  // 社交信息删除
  socialInfoDelete: (req: BlogSocialInfoDeleteReq) =>
    blogSocialInfoDelete(req)
}
