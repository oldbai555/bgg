package conf

import (
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/bgg/service/lbblog"
)

var Global = webtool.GenWebToolByYaml(lbblog.ServerName, webtool.OptionWithOrm(), webtool.OptionWithRdb(), webtool.OptionWithServer())
