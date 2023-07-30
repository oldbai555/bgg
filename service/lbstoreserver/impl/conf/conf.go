package conf

import (
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbstore"
)

var Global = webtool.GenWebToolByYaml(lbstore.ServerName, webtool.OptionWithOrm(), webtool.OptionWithRdb(), webtool.OptionWithServer())
