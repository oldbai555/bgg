import {
  publicBlogArticleList,
  publicBlogArticleDetail,
  publicBlogFriendLinkList,
  publicBlogSocialInfoList,
  publicBlogTagList,
  publicBlogAuthorInfo,
  publicBlogArticleStats,
  publicBlogArticlePrev,
  publicBlogArticleNext
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
  PublicBlogArticleNextResp
} from '@/api/generated/admin'

/**
 * 博客相关 API 封装
 */
export const blogApi = {
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
    publicBlogArticleNext(req) as Promise<PublicBlogArticleNextResp>
}
