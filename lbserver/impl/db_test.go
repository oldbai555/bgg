package impl

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/lbconst"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/bgg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"testing"
)

func init() {
	lb = &Tool{}
	lb.WebTool, _ = webtool.NewWebTool(webtool.OptionWithOrm(&lbuser.ModelUser{}), webtool.OptionWithRdb())
	InitDbOrm()
}

func TestGenUserOrm(t *testing.T) {
	var user lbuser.ModelUser
	err := UserOrm.NewScope().UpdateOrCreate(context.Background(), map[string]interface{}{
		"username": "oldbai",
		"password": utils.StrMd5("123456"),
	}, map[string]interface{}{
		"password": utils.StrMd5("123456"),
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
