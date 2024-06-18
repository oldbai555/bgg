package lb

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/gormx"
	"reflect"
)

type Scope struct {
	*gormx.Scope

	size      uint32
	page      uint32
	skipTotal bool
}

type Model struct {
	*gormx.Model
}

func NewModel(db *gorm.DB, m gorm.Tabler, err error) *Model {
	return &Model{
		Model: gormx.NewModel(db, m, err),
	}
}

func (f *Model) NewScope(ctx context.Context) *Scope {
	return &Scope{
		Scope: f.Model.NewScope(ctx),
	}
}

func (f *Model) NewList(ctx context.Context, listOption *ListOption) *Scope {
	return &Scope{
		Scope: f.Model.NewScope(ctx),

		size:      listOption.GetSize(),
		page:      listOption.GetPage(),
		skipTotal: listOption.GetSkipTotal(),
	}
}

func (p *Scope) Corp(corpId uint32) *Scope {
	p.Eq("corp_id", corpId)
	return p
}

// FindPaginate 分页查找
func (p *Scope) FindPaginate(list interface{}) (*Paginate, error) {
	var total int64
	if !p.skipTotal {
		err := p.DB().Count(&total).Error
		if err != nil {
			return nil, err
		}
	}

	var page = uint32(0)
	if p.page-1 > 0 {
		page = p.page - 1
	}

	err := p.DB().Limit(int(p.size)).Offset(int(page * p.size)).Find(list).Error
	if err != nil {
		return nil, err
	}

	return &Paginate{
		Total: uint64(total),
		Size:  p.size,
		Page:  p.page,
	}, nil
}

func (p *Scope) Chunk(limit int, list interface{}, cb func() (stop bool)) error {
	if cb == nil {
		return nil
	}

	var total int64
	err := p.DB().Count(&total).Error
	if err != nil {
		return err
	}
	maxPage := total / int64(limit)
	if total%int64(limit) != 0 {
		maxPage += 1
	}

	var page = 0
	ClearSlice(list)
	for maxPage > int64(page) {
		ClearSlice(list)
		err = p.DB().Limit(limit).Offset(page * limit).Find(list).Error
		if err != nil {
			return err
		}
		if cb() {
			return nil
		}
		page++
	}
	return nil
}

func ClearSlice(ptr interface{}) {
	if ptr == nil {
		return
	}
	vo := reflect.ValueOf(ptr)
	if vo.Kind() != reflect.Ptr {
		panic("required ptr to slice type")
	}
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	if vo.Kind() != reflect.Slice {
		panic("required ptr to slice type")
	}
	vo.Set(reflect.MakeSlice(vo.Type(), 0, 0))
}

func (f *Model) FirstOrCreate(ctx context.Context, candMap map[string]interface{}, out interface{}) error {
	err := f.NewScope(ctx).AndMap(candMap).First(out)
	if err != nil && !f.IsNotFoundErr(err) {
		log.Errorf("err:%v", err)
		return err
	}

	optDb := f.NewScope(ctx)
	if f.IsNotFoundErr(err) {
		err = mapstructure.Decode(candMap, out)
		if err != nil {
			return err
		}

		_, err := optDb.Create(out)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
