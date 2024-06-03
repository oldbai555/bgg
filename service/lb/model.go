package lb

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/gormx"
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
			log.Errorf("err is %v", err)
			return nil, err
		}
	}

	var page = uint32(0)
	if p.page-1 > 0 {
		page = p.page - 1
	}

	err := p.DB().Limit(int(p.size)).Offset(int(page * p.size)).Find(list).Error
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	return &Paginate{
		Total: uint64(total),
		Size:  p.size,
		Page:  p.page,
	}, nil
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
			log.Errorf("err is : %v", err)
			return err
		}

		_, err := optDb.Create(out)
		if err != nil {
			log.Errorf("err is %v", err)
			return err
		}
		return nil
	}

	return nil
}
