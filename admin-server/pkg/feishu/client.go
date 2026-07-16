// Package feishu 封装飞书开放平台（自建应用）网页登录所需的最小 OpenAPI 调用：
// 用授权 code 换 user access_token，再用 access_token 换用户身份信息。
// 端点/字段核实自官方文档（2026-07）：
//   - https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/authentication-management/access-token/get-user-access-token
//   - https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/user_info/get
package feishu

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"postapocgame/admin-server/pkg/errs"
)

const (
	userAccessTokenURL = "https://open.feishu.cn/open-apis/authen/v2/oauth/token"
	userInfoURL        = "https://open.feishu.cn/open-apis/authen/v1/user_info"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// Client 飞书自建应用客户端，AppId/AppSecret 来自 FeishuConf。
type Client struct {
	AppId       string
	AppSecret   string
	RedirectUri string
}

func NewClient(appId, appSecret, redirectUri string) *Client {
	return &Client{AppId: appId, AppSecret: appSecret, RedirectUri: redirectUri}
}

// UserInfo 授权 code 换回的用户身份信息。
type UserInfo struct {
	OpenId    string
	UnionId   string
	Name      string
	AvatarUrl string
	Mobile    string
	Email     string
}

type userAccessTokenResp struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	AccessToken string `json:"access_token"`
}

type userInfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OpenId    string `json:"open_id"`
		UnionId   string `json:"union_id"`
		Name      string `json:"name"`
		AvatarUrl string `json:"avatar_url"`
		Mobile    string `json:"mobile"`
		Email     string `json:"email"`
	} `json:"data"`
}

// ExchangeUserInfo 用授权 code 换取飞书用户身份信息：先用 code 换 user access_token，
// 再用 access_token 查用户信息。
func (c *Client) ExchangeUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	accessToken, err := c.fetchUserAccessToken(ctx, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "创建飞书用户信息请求失败", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadGateway, "请求飞书用户信息接口失败", err)
	}
	defer resp.Body.Close()

	var result userInfoResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.Wrap(errs.CodeBadGateway, "解析飞书用户信息响应失败", err)
	}
	if result.Code != 0 {
		return nil, errs.New(errs.CodeBadGateway, "查询飞书用户信息失败: "+result.Msg)
	}

	return &UserInfo{
		OpenId:    result.Data.OpenId,
		UnionId:   result.Data.UnionId,
		Name:      result.Data.Name,
		AvatarUrl: result.Data.AvatarUrl,
		Mobile:    result.Data.Mobile,
		Email:     result.Data.Email,
	}, nil
}

func (c *Client) fetchUserAccessToken(ctx context.Context, code string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     c.AppId,
		"client_secret": c.AppSecret,
		"code":          code,
		"redirect_uri":  c.RedirectUri,
	})
	if err != nil {
		return "", errs.Wrap(errs.CodeInternalError, "构造飞书 access_token 请求失败", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userAccessTokenURL, strings.NewReader(string(body)))
	if err != nil {
		return "", errs.Wrap(errs.CodeInternalError, "创建飞书 access_token 请求失败", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", errs.Wrap(errs.CodeBadGateway, "请求飞书 access_token 接口失败", err)
	}
	defer resp.Body.Close()

	var result userAccessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errs.Wrap(errs.CodeBadGateway, "解析飞书 access_token 响应失败", err)
	}
	if result.Code != 0 {
		return "", errs.New(errs.CodeBadGateway, "获取飞书 access_token 失败: "+result.Msg)
	}

	return result.AccessToken, nil
}
