package webtool

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"reflect"
	"strings"
)

type OrmCondBuilder struct {
	db *gorm.DB

	limit     uint32
	offset    uint32
	skipCount bool
}

func NewCondBuilder(db *gorm.DB) *OrmCondBuilder {
	return &OrmCondBuilder{
		db: db,
	}
}

func NewList(db *gorm.DB, listOption *lbconst.ListOption) *OrmCondBuilder {
	return &OrmCondBuilder{
		db: db,

		limit:     listOption.Limit,
		offset:    listOption.Offset,
		skipCount: listOption.SkipCount,
	}
}

func (p *OrmCondBuilder) Eq(field string, v interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("`%s` = ?", field), v)
	return p
}

func (p *OrmCondBuilder) NotEq(f string, v interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("`%s` != ?", f), v)
	return p
}

// AndMap 示例
// key: "name" val: "hangman" sql = `name` = "hangman"
// key: "name like" val: %hangman% sql = `name` like "%hangman%"
func (p *OrmCondBuilder) AndMap(kv map[string]interface{}) *OrmCondBuilder {
	if len(kv) > 0 {
		var condList []string
		var argList []interface{}
		for k, v := range kv {
			if k == "" {
				panic(any("invalid empty key"))
			}
			split := strings.Split(k, " ")
			if len(split) == 2 {
				condList = append(condList, fmt.Sprintf("(`%s` %s ?)", split[0], split[1]))
			} else {
				condList = append(condList, fmt.Sprintf("`%s` = ?", k))
			}
			argList = append(argList, v)
		}
		cond := strings.Join(condList, " AND ")
		p.db.Where(cond, argList...)
	}
	return p
}

// OrMap 示例
// key: "name" val: "hangman" sql = `name` = "hangman"
// key: "name like" val: %hangman% sql = `name` like "%hangman%"
func (p *OrmCondBuilder) OrMap(kv map[string]interface{}) *OrmCondBuilder {
	if len(kv) > 0 {
		var condList []string
		var argList []interface{}
		for k, v := range kv {
			if k == "" {
				panic(any("invalid empty key"))
			}
			split := strings.Split(k, " ")
			if len(split) == 2 {
				condList = append(condList, fmt.Sprintf("(`%s` %s ?)", split[0], split[1]))
			} else {
				condList = append(condList, fmt.Sprintf("(`%s` = ?)", k))
			}
			argList = append(argList, v)
		}
		cond := strings.Join(condList, " OR ")
		p.db.Where(cond, argList...)
	}
	return p
}

func (p *OrmCondBuilder) Like(f, v string) *OrmCondBuilder {
	if v != "" {
		v := utils.EscapeMysqlLikeWildcardIgnore2End(v)
		p.db.Where(fmt.Sprintf("`%s` LIKE ?", f), v)
	}
	return p
}

func (p *OrmCondBuilder) NotLike(f, v string) *OrmCondBuilder {
	if v != "" {
		v = utils.EscapeMysqlLikeWildcardIgnore2End(v)
		v = utils.QuoteName(fmt.Sprintf("%%%s%%", v))
		p.db.Where(
			fmt.Sprintf("`%s` NOT LIKE %s", f, v))
	}
	return p
}

func (p *OrmCondBuilder) In(f string, i interface{}) *OrmCondBuilder {
	v := reflect.ValueOf(i)
	if v.Type().Kind() != reflect.Slice {
		panic(any("invalid input type, slice"))
	}
	if v.Len() == 0 {
		p.db.Where("1=0")
		return p
	}
	p.db.Where(fmt.Sprintf("`%s` in (?)", f), utils.UniqueSliceV2(i))
	return p
}

func (p *OrmCondBuilder) NotIn(f string, i interface{}) *OrmCondBuilder {
	v := reflect.ValueOf(i)
	// 如果不是slice，也是可以的，比如 id in (1)
	if v.Type().Kind() != reflect.Slice {
		panic(any("invalid input type, slice"))
	}
	if v.Len() == 0 {
		p.db.Where("1=0")
		return p
	}
	p.db.Where(fmt.Sprintf("`%s` not in (?)", f), utils.UniqueSliceV2(i))
	return p
}

func (p *OrmCondBuilder) Lt(f string, v interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("`%s` < ?", f), v)
	return p
}

func (p *OrmCondBuilder) Lte(f string, v interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("`%s` <= ?", f), v)
	return p
}

func (p *OrmCondBuilder) Gt(f string, v interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("`%s` > ?", f), v)
	return p
}

func (p *OrmCondBuilder) Gte(f string, v interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("`%s` >= ?", f), v)
	return p
}

func (p *OrmCondBuilder) Order(order string) *OrmCondBuilder {
	p.db.Order(order)
	return p
}

func (p *OrmCondBuilder) OrderByDesc(order ...string) *OrmCondBuilder {
	p.db.Order(fmt.Sprintf("%s DESC", strings.Join(order, ",")))
	return p
}

func (p *OrmCondBuilder) OrderByAsc(order ...string) *OrmCondBuilder {
	p.db.Order(fmt.Sprintf("%s ASC", strings.Join(order, ",")))
	return p
}

func (p *OrmCondBuilder) Group(group ...string) *OrmCondBuilder {
	for _, s := range group {
		p.db.Group(s)
	}
	return p
}

// Between 相当于 field >= min || field <= max
func (p *OrmCondBuilder) Between(fieldName string, min, max interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", utils.QuoteFieldName(fieldName)), min, max)
	return p
}

// NotBetween 相当于 field < min || field > max
func (p *OrmCondBuilder) NotBetween(fieldName string, min, max interface{}) *OrmCondBuilder {
	p.db.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", utils.QuoteFieldName(fieldName)), min, max)
	return p
}

// UnScoped 去除逻辑删除条件
func (p *OrmCondBuilder) UnScoped() {
	p.db.Unscoped()
}

// First 查找
func (p *OrmCondBuilder) First(ctx context.Context, dest interface{}) error {
	return p.db.First(dest).Error
}

// Update 更新
func (p *OrmCondBuilder) Update(ctx context.Context, values map[string]interface{}) (int64, error) {
	res := p.db.Updates(values)
	if res.Error != nil {
		log.Errorf("err is %v", res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// Delete 删除
func (p *OrmCondBuilder) Delete(ctx context.Context, obj interface{}) (int64, error) {
	res := p.db.Delete(obj)
	if res.Error != nil {
		log.Errorf("err is %v", res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// Create 创建
func (p *OrmCondBuilder) Create(ctx context.Context, dest interface{}) (int64, error) {
	res := p.db.Create(dest)
	if res.Error != nil {
		log.Errorf("err is %v", res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// Find 查找所有
func (p *OrmCondBuilder) Find(ctx context.Context, dest interface{}) error {
	err := p.db.Find(dest).Error
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

// FindPage 分页查找
func (p *OrmCondBuilder) FindPage(ctx context.Context, list interface{}) (*lbconst.Page, error) {
	var total int64
	if !p.skipCount {
		err := p.db.Count(&total).Error
		if err != nil {
			log.Errorf("err is %v", err)
			return nil, err
		}
	}
	err := p.db.Limit(int(p.limit)).Offset(int(p.offset)).Find(list).Error
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &lbconst.Page{
		Total:  uint64(total),
		Limit:  p.limit,
		Offset: p.offset,
	}, nil
}
