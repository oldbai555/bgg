package conf

import (
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbuser"
)

var Global = webtool.GenWebToolByYaml(lbuser.ServerName, webtool.OptionWithOrm(), webtool.OptionWithRdb(), webtool.OptionWithServer())
