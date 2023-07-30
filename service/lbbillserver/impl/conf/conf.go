package conf

import (
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbbill"
)

var Global = webtool.GenWebToolByYaml(lbbill.ServerName, webtool.OptionWithOrm(), webtool.OptionWithRdb(), webtool.OptionWithServer())
