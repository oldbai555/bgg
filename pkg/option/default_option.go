package option

import (
	"github.com/oldbai555/bgg/internal/bgorm"
	"github.com/oldbai555/bgg/service/lb"
	"github.com/oldbai555/lbtool/log"
)

func ProcessDefaultOptions(listOption *lb.Options, db *bgorm.Scope) error {
	err := lb.NewOptionsProcessor(listOption).
		AddStringList(
			lb.DefaultListOption_DefaultListOptionSelect,
			func(valList []string) error {
				db.Select(valList)
				return nil
			}).
		AddUint32(
			lb.DefaultListOption_DefaultListOptionOrderBy,
			func(val uint32) error {
				if val == uint32(lb.DefaultOrderBy_DefaultOrderByCreatedDesc) {
					db.OrderByDesc("created_at")
				} else if val == uint32(lb.DefaultOrderBy_DefaultOrderByCreatedAcs) {
					db.OrderByAsc("created_at")
				} else if val == uint32(lb.DefaultOrderBy_DefaultOrderByIdDesc) {
					db.OrderByDesc("id")
				}
				return nil
			}).
		AddStringList(
			lb.DefaultListOption_DefaultListOptionGroupBy,
			func(valList []string) error {
				db.Group(valList...)
				return nil
			}).
		AddBool(
			lb.DefaultListOption_DefaultListOptionWithTrash,
			func(val bool) error {
				if val {
					db.UnScoped()
				}
				return nil
			}).
		AddUint64List(
			lb.DefaultListOption_DefaultListOptionIdList,
			func(valList []uint64) error {
				if len(valList) == 1 {
					db.Eq("id", valList[0])
				} else {
					db.In("id", valList)
				}
				return nil
			}).
		AddTimeStampRange(
			lb.DefaultListOption_DefaultListOptionCreatedAt,
			func(begin, end uint32) error {
				db.Between("created_at", begin, end)
				return nil
			}).
		AddUint32List(
			lb.DefaultListOption_DefaultListOptionCreatorIdList,
			func(valList []uint32) error {
				db.In("creator_id", valList)
				return nil
			}).
		//AddUint64List(
		//	lb.DefaultListOption_DefaultListOptionCorpIdList,
		//	func(valList []uint64) error {
		//		if len(valList) == 1 {
		//			db.Eq("corp_id", valList[0])
		//		} else {
		//			db.In("corp_id", valList)
		//		}
		//		return nil
		//	}).
		Process()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
