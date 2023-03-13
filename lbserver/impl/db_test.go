package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/bgg/client/lbuser"
	webtool2 "github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"testing"
)

func init() {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	lb = &Tool{}
	lb.WebTool, _ = webtool2.NewWebTool(v, webtool2.OptionWithOrm(&lbuser.ModelUser{}), webtool2.OptionWithRdb(), webtool2.OptionWithStorage())
	InitDbOrm()
}

func TestGenUserOrm(t *testing.T) {
	var user lbuser.ModelUser
	err := UserOrm.NewScope().UpdateOrCreate(context.Background(), map[string]interface{}{
		"username": "superadmin",
		"password": utils.StrMd5("123456"),
	}, map[string]interface{}{
		"password": utils.StrMd5("123456"),
		"role":     1,
	}, &user)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	log.Infof("user is %v", user)
}

func TestFindUser(t *testing.T) {
	var users []*lbuser.ModelUser
	page, err := UserOrm.NewList(lbconst.NewListOption()).Like(lbuser.FieldUsername_, fmt.Sprintf("%%%s%%", "oldbai")).FindPage(context.TODO(), &users)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	log.Infof("users is %v", users)
	log.Infof("page is %v", page)
}

func TestGenToken(t *testing.T) {
	log.Infof("token is %v", utils.StrMd5("Aroot@123"))
}
