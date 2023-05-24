package webtool

import (
	"github.com/oldbai555/bgg/client/lb"
	"github.com/oldbai555/lbtool/log"
)

func ProcessDefaultOptions(options *lb.Options, db *OrmCondBuilder) error {
	err := lb.NewOptionsProcessor(options).
		AddUint32(lb.DefaultOption_DefaultOptionOrderBy, func(val uint32) error {
			switch lb.DefaultOrderBy(val) {
			case lb.DefaultOrderBy_DefaultOrderByCreatedAtAsc:
				db.OrderByAsc("created_at")
			case lb.DefaultOrderBy_DefaultOrderByCreatedAtDesc:
				db.OrderByDesc("created_at")

			}
			return nil
		}).
		AddTimeStampRange(lb.DefaultOption_DefaultOptionCreatedAtRange, func(beginAt, endAt uint32) error {
			db.Between("created_at", beginAt, endAt)
			return nil
		}).
		Process()
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}
