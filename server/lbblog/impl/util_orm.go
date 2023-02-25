package impl

import (
	"context"
	"fmt"
	"github.com/huandu/go-clone"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/bgg/client/lbconst"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"reflect"
	"strings"
)

type ResetDbFunc func() *gorm.DB

// OrmCondBuilder 扩展条件构造器
type OrmCondBuilder struct {
	obj interface{}
	err error
}

type Scope struct {
	db *gorm.DB
	//condList pie.Strings
	//orders   pie.Strings
	//group    pie.Strings

	obj interface{}
	err error

	limit     uint32
	offset    uint32
	skipCount bool
}

func NewOrmCondBuilder(obj interface{}, err error) *OrmCondBuilder {
	return &OrmCondBuilder{
		obj: obj,
		err: err,
	}
}

func (p *OrmCondBuilder) NewScope() *Scope {
	newObj := clone.Clone(p.obj)
	newErr := clone.Clone(p.err).(*lberr.LbErr)
	return &Scope{
		db:  lb.Orm.Model(newObj),
		err: newErr,
		obj: newObj,
	}
}

func (p *OrmCondBuilder) NewList(listOption *lbconst.ListOption) *Scope {
	return &Scope{
		db:        lb.Orm.Model(p.obj),
		err:       p.err,
		obj:       p.obj,
		limit:     listOption.Limit,
		offset:    listOption.Offset,
		skipCount: listOption.SkipCount,
	}
}

func (p *OrmCondBuilder) IsNotFoundErr(err error) bool {
	return err == p.err
}

func (p *Scope) CopyScope() *Scope {
	return &Scope{
		db:  lb.Orm.Model(p.obj),
		err: p.err,
		obj: p.obj,
	}
}

func (p *Scope) Eq(field string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` = ?", field), v)
	return p
}

func (p *Scope) NotEq(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` != ?", f), v)
	return p
}

// AndMap 示例
// key: name val: zhangsan sql = `name` = zhangsan
// key: name like val: %zhangsan% sql = `name` like %zhangsan%
func (p *Scope) AndMap(kv map[string]interface{}) *Scope {
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
// key: name val: zhangsan sql = `name` = zhangsan
// key: name like val: %zhangsan% sql = `name` like %zhangsan%
func (p *Scope) OrMap(kv map[string]interface{}) *Scope {
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

func (p *Scope) Like(f, v string) *Scope {
	if v != "" {
		v := utils.EscapeMysqlLikeWildcardIgnore2End(v)
		p.db.Where(fmt.Sprintf("`%s` LIKE ?", f), v)
	}
	return p
}

func (p *Scope) NotLike(f, v string) *Scope {
	if v != "" {
		v = utils.EscapeMysqlLikeWildcardIgnore2End(v)
		v = utils.QuoteName(fmt.Sprintf("%%%s%%", v))
		p.db.Where(
			fmt.Sprintf("`%s` NOT LIKE %s", f, v))
	}
	return p
}

func (p *Scope) In(f string, i interface{}) *Scope {
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

func (p *Scope) NotIn(f string, i interface{}) *Scope {
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

func (p *Scope) Lt(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` < ?", f), v)
	return p
}

func (p *Scope) Lte(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` <= ?", f), v)
	return p
}

func (p *Scope) Gt(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` > ?", f), v)
	return p
}

func (p *Scope) Gte(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` >= ?", f), v)
	return p
}

func (p *Scope) Order(order string) *Scope {
	p.db.Order(order)
	return p
}

func (p *Scope) OrderByDesc(order ...string) *Scope {
	p.db.Order(fmt.Sprintf("%s DESC", strings.Join(order, ",")))
	return p
}

func (p *Scope) OrderByAsc(order ...string) *Scope {
	p.db.Order(fmt.Sprintf("%s ASC", strings.Join(order, ",")))
	return p
}

func (p *Scope) Group(group ...string) *Scope {
	for _, s := range group {
		p.db.Group(s)
	}
	return p
}

// Between 相当于 field >= min || field <= max
func (p *Scope) Between(fieldName string, min, max interface{}) *Scope {
	p.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", utils.QuoteFieldName(fieldName)), min, max)
	return p
}

// NotBetween 相当于 field < min || field > max
func (p *Scope) NotBetween(fieldName string, min, max interface{}) *Scope {
	p.db.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", utils.QuoteFieldName(fieldName)), min, max)
	return p
}

func (p *Scope) UnScoped() {
	p.db.Unscoped()
}

// First 查找
func (p *Scope) First(ctx context.Context, dest interface{}) error {
	err := p.db.First(dest).Error
	if err == gorm.ErrRecordNotFound {
		return p.err
	}
	return err
}

// IsEmpty 是否存在
func (p *Scope) IsEmpty(ctx context.Context) (bool, error) {
	dest := clone.Clone(p.obj)
	err := p.db.First(dest).Error
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	return gorm.ErrRecordNotFound == err, nil
}

// Update 更新
func (p *Scope) Update(ctx context.Context, values map[string]interface{}) error {
	err := p.db.Updates(values).Error
	return err
}

// Delete 删除
func (p *Scope) Delete(ctx context.Context) error {
	err := p.db.Delete(p.obj).Error
	return err
}

// Create 创建
func (p *Scope) Create(ctx context.Context, dest interface{}) error {
	err := p.db.Create(dest).Error
	return err
}

// UpdateOrCreate 更新或创建
// cond 条件
// attr 更新的属性
func (p *Scope) UpdateOrCreate(ctx context.Context, cond map[string]interface{}, attr map[string]interface{}, dest interface{}) error {
	err := p.CopyScope().AndMap(cond).First(ctx, dest)
	if p.err == err {
		for k, v := range cond {
			attr[k] = v
		}
		err = mapstructure.Decode(attr, dest)
		if err != nil {
			log.Errorf("err is : %v", err)
			return err
		}
		err = p.CopyScope().Create(ctx, dest)
		if err != nil {
			log.Errorf("err is %v", err)
			return err
		}
		return nil
	}
	if err != nil && p.err != err {
		return err
	}
	return p.CopyScope().AndMap(cond).Update(ctx, attr)
}

func (p *Scope) Find(ctx context.Context, dest interface{}) error {
	err := p.db.Find(dest).Error
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func (p *Scope) FindPage(ctx context.Context, list interface{}) (*lbconst.Page, error) {
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
