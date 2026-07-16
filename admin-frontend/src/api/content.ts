import {
  // 公开博客接口
  publicBlogArticleList,
  publicBlogArticleDetail,
  publicBlogFriendLinkList,
  publicBlogSocialInfoList,
  publicBlogTagList,
  publicBlogAuthorInfo,
  publicBlogArticleStats,
  publicBlogArticlePrev,
  publicBlogArticleNext,
  // 后台博客文章
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
  // 后台博客标签
  blogTagList,
  blogTagCreate,
  blogTagUpdate,
  blogTagDelete,
  blogTagOptions,
  // 后台博客友情链接
  blogFriendLinkList,
  blogFriendLinkCreate,
  blogFriendLinkUpdate,
  blogFriendLinkDelete,
  // 后台博客社交信息
  blogSocialInfoList,
  blogSocialInfoCreate,
  blogSocialInfoUpdate,
  blogSocialInfoDelete,
  // 公开视频接口
  publicVideoList,
  publicVideoDetail,
  // 后台视频管理
  videoList,
  videoCreate,
  videoUpdate,
  videoDelete,
  videoCollect,
  videoCollectOptions,
  // m3u8 代理
  m3u8Proxy
} from '@/api/generated/admin'
import type {
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
  BlogArticleListReq,
  BlogArticleListResp,
  BlogArticleDetailReq,
  BlogArticleDetailResp,
  BlogTagListReq,
  BlogTagListResp,
  BlogTagOptionsReq,
  BlogTagOptionsResp,
  BlogFriendLinkListReq,
  BlogFriendLinkListResp,
  BlogSocialInfoListReq,
  BlogSocialInfoListResp,
  PublicVideoListReq,
  PublicVideoListResp,
  PublicVideoDetailReq,
  PublicVideoDetailResp
} from '@/api/generated/admin'

/**
 * Content 域 API 封装（博客 + 视频，对应后端已合并的 content-rpc）
 * 原 src/api/blog.ts + src/api/video.ts 合并到此处，按原分组保留注释分节。
 */
export const contentApi = {
  // ========== 博客：公开接口 ==========
  publicArticleList: (req: PublicBlogArticleListReq) =>
    publicBlogArticleList(req) as Promise<PublicBlogArticleListResp>,
  publicArticleDetail: (req: PublicBlogArticleDetailReq) =>
    publicBlogArticleDetail(req) as Promise<PublicBlogArticleDetailResp>,
  publicFriendLinkList: () => publicBlogFriendLinkList() as Promise<PublicBlogFriendLinkListResp>,
  publicSocialInfoList: () => publicBlogSocialInfoList() as Promise<PublicBlogSocialInfoListResp>,
  publicTagList: () => publicBlogTagList() as Promise<PublicBlogTagListResp>,
  publicAuthorInfo: () => publicBlogAuthorInfo() as Promise<PublicBlogAuthorInfoResp>,
  publicArticleStats: () => publicBlogArticleStats() as Promise<PublicBlogArticleStatsResp>,
  publicArticlePrev: (req: PublicBlogArticlePrevReq) =>
    publicBlogArticlePrev(req) as Promise<PublicBlogArticlePrevResp>,
  publicArticleNext: (req: PublicBlogArticleNextReq) =>
    publicBlogArticleNext(req) as Promise<PublicBlogArticleNextResp>,

  // ========== 博客：后台文章 ==========
  articleList: (req: BlogArticleListReq) => blogArticleList(req) as Promise<BlogArticleListResp>,
  articleCreate: blogArticleCreate,
  articleUpdate: blogArticleUpdate,
  articleDelete: blogArticleDelete,
  articleDetail: (req: BlogArticleDetailReq) => blogArticleDetail(req) as Promise<BlogArticleDetailResp>,
  articlePublish: blogArticlePublish,
  articleSubmit: blogArticleSubmit,
  articleTop: blogArticleTop,
  articleUnpublish: blogArticleUnpublish,
  articleUntop: blogArticleUntop,
  articleAudit: blogArticleAudit,
  articleAuditUnpublish: blogArticleAuditUnpublish,

  // ========== 博客：标签 ==========
  tagList: (req: BlogTagListReq) => blogTagList(req) as Promise<BlogTagListResp>,
  tagCreate: blogTagCreate,
  tagUpdate: blogTagUpdate,
  tagDelete: blogTagDelete,
  tagOptions: (req: BlogTagOptionsReq) => blogTagOptions(req) as Promise<BlogTagOptionsResp>,

  // ========== 博客：友情链接 ==========
  friendLinkList: (req: BlogFriendLinkListReq) =>
    blogFriendLinkList(req) as Promise<BlogFriendLinkListResp>,
  friendLinkCreate: blogFriendLinkCreate,
  friendLinkUpdate: blogFriendLinkUpdate,
  friendLinkDelete: blogFriendLinkDelete,

  // ========== 博客：社交信息 ==========
  socialInfoList: (req: BlogSocialInfoListReq) =>
    blogSocialInfoList(req) as Promise<BlogSocialInfoListResp>,
  socialInfoCreate: blogSocialInfoCreate,
  socialInfoUpdate: blogSocialInfoUpdate,
  socialInfoDelete: blogSocialInfoDelete,

  // ========== 视频：公开接口 ==========
  publicVideoList: (req: PublicVideoListReq) => publicVideoList(req) as Promise<PublicVideoListResp>,
  publicVideoDetail: (req: PublicVideoDetailReq) =>
    publicVideoDetail(req) as Promise<PublicVideoDetailResp>,

  // ========== 视频：后台管理 ==========
  videoList,
  videoCreate,
  videoUpdate,
  videoDelete,
  videoCollect,
  videoCollectOptions,

  // ========== 视频：m3u8 代理 ==========
  m3u8Proxy
}
