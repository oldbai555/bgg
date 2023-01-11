package impl

import (
	"context"
	"github.com/storyicon/grbac"
)

func LoadAuthorizationRules() (rules grbac.Rules, err error) {
	// ID=0 的规则表明任何具有任何角色的人都可以访问所有资源。
	// 但是ID=1的规则表明只有 lb_user 可以对用户进行增删改操作。
	rules = grbac.Rules{
		{
			ID: 0,
			Resource: &grbac.Resource{
				Host:   "*",
				Path:   "*",
				Method: "*",
			},
			Permission: &grbac.Permission{
				AllowAnyone:     false,
				AuthorizedRoles: []string{"*"},
				ForbiddenRoles:  []string{},
			},
		},
		{
			ID: 1,
			Resource: &grbac.Resource{
				Host:   "*",
				Path:   "/user/**",
				Method: "{POST,GET}",
			},
			Permission: &grbac.Permission{
				AuthorizedRoles: []string{"lb_user"},
				ForbiddenRoles:  []string{},
				AllowAnyone:     false,
			},
		},
		{
			ID: 1,
			Resource: &grbac.Resource{
				Host:   "*",
				Path:   "/blog/**",
				Method: "{POST,GET}",
			},
			Permission: &grbac.Permission{
				AuthorizedRoles: []string{"lb_user"},
				ForbiddenRoles:  []string{},
				AllowAnyone:     false,
			},
		},
	}
	return
}

func GetUserRoles(ctx context.Context, uid uint64) ([]string, error) {
	return []string{"lb_user"}, nil
}
