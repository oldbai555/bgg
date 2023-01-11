package lbconst

import (
	"fmt"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	ListOptionValueTypeNil            = 0
	ListOptionValueTypeUint32         = 1
	ListOptionValueTypeString         = 2
	ListOptionValueTypeUint32List     = 3
	ListOptionValueTypeUint64         = 4
	ListOptionValueTypeTimeStampRange = 5
	ListOptionValueTypeUint64List     = 6
	ListOptionValueTypeBool           = 7
	ListOptionValueTypeStringList     = 8
	ListOptionValueTypeInt32          = 9
)

const (
	defaultLimit = 2000
)

type ListOptionHandler struct {
	typ              int
	cbNone           func() error
	cbUint32         func(val uint32) error
	cbInt32          func(val int32) error
	cbString         func(val string) error
	cbUint32List     func(valList []uint32) error
	cbUint64         func(val uint64) error
	cbTimeStampRange func(beginAt, endAt uint32) error
	cbUint64List     func(val []uint64) error
	cbBool           func(val bool) error
	cbStringList     func(val []string) error
	ignoreZeroValue  bool
}
type ListOptionProcessor struct {
	listOption *ListOption
	handlers   map[int32]*ListOptionHandler
}

func toInt32(i interface{}) int32 {
	t := reflect.TypeOf(i)
	k := t.Kind()
	switch k {
	case reflect.Int,
		reflect.Int32,
		reflect.Int64,
		reflect.Int16,
		reflect.Int8:
		return int32(reflect.ValueOf(i).Int())
	case reflect.Uint,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint16,
		reflect.Uint8:
		return int32(reflect.ValueOf(i).Uint())
	}
	return 0
}

func NewListOptionProcessor(listOption *ListOption) *ListOptionProcessor {
	return &ListOptionProcessor{
		listOption: listOption,
		handlers:   make(map[int32]*ListOptionHandler),
	}
}

func (p *ListOptionProcessor) AddNone(typ interface{}, cb func() error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:    ListOptionValueTypeNil,
		cbNone: cb,
	}
	return p
}

func (p *ListOptionProcessor) AddUint32(typ interface{}, cb func(val uint32) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:      ListOptionValueTypeUint32,
		cbUint32: cb,
	}
	return p
}

func (p *ListOptionProcessor) AddInt32(typ interface{}, cb func(val int32) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:     ListOptionValueTypeInt32,
		cbInt32: cb,
	}
	return p
}

func (p *ListOptionProcessor) AddString(typ interface{}, cb func(val string) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:             ListOptionValueTypeString,
		cbString:        cb,
		ignoreZeroValue: true,
	}
	return p
}

func (p *ListOptionProcessor) AddStringIgnoreZero(typ interface{}, cb func(val string) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:      ListOptionValueTypeString,
		cbString: cb,
	}
	return p
}

func (p *ListOptionProcessor) AddStringList(typ interface{}, cb func(valList []string) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:             ListOptionValueTypeStringList,
		cbStringList:    cb,
		ignoreZeroValue: true,
	}
	return p
}

func (p *ListOptionProcessor) AddUint32List(typ interface{}, cb func(valList []uint32) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:             ListOptionValueTypeUint32List,
		cbUint32List:    cb,
		ignoreZeroValue: true,
	}
	return p
}

func (p *ListOptionProcessor) AddUint64(typ interface{}, cb func(val uint64) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:      ListOptionValueTypeUint64,
		cbUint64: cb,
	}
	return p
}

func (p *ListOptionProcessor) AddUint64List(typ interface{}, cb func(valList []uint64) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:             ListOptionValueTypeUint64List,
		cbUint64List:    cb,
		ignoreZeroValue: true,
	}
	return p
}

func (p *ListOptionProcessor) AddTimeStampRange(typ interface{}, cb func(beginAt, endAt uint32) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:              ListOptionValueTypeTimeStampRange,
		cbTimeStampRange: cb,
		ignoreZeroValue:  true,
	}
	return p
}

func (p *ListOptionProcessor) AddBool(typ interface{}, cb func(val bool) error) *ListOptionProcessor {
	x := toInt32(typ)
	p.handlers[x] = &ListOptionHandler{
		typ:    ListOptionValueTypeBool,
		cbBool: cb,
	}
	return p
}

func (p *ListOptionProcessor) Process() error {
	if p.listOption == nil || p.handlers == nil || len(p.handlers) == 0 {
		return nil
	}
	var err error
	for _, v := range p.listOption.Options {
		h := p.handlers[v.Type]
		if h == nil {
			continue
		}
		switch h.typ {
		case ListOptionValueTypeNil:
			if h.cbNone != nil {
				err = h.cbNone()
				if err != nil {
					return err
				}
			}

		case ListOptionValueTypeUint32:
			if v.Value == "" {
				continue
			}
			x, err := strconv.ParseInt(v.Value, 10, 32)
			if err != nil {
				return lberr.NewInvalidArg("invalid option value with type %d, expected uint32", v.Type)
			}
			if h.cbUint32 != nil {
				err = h.cbUint32(uint32(x))
				if err != nil {
					return err
				}
			}

		case ListOptionValueTypeString:
			if v.Value == "" && h.ignoreZeroValue {
				continue
			}
			if h.cbString != nil {
				if err = h.cbString(v.Value); err != nil {
					return err
				}
			}

		case ListOptionValueTypeStringList:
			if v.Value == "" && h.ignoreZeroValue {
				continue
			}
			if h.cbStringList != nil {
				list := strings.Split(v.Value, ",")
				// 过滤掉空串
				var nonEmptyList []string
				for _, v := range list {
					if v != "" {
						nonEmptyList = append(nonEmptyList, v)
					}
				}
				if len(nonEmptyList) == 0 && h.ignoreZeroValue {
					continue
				}
				if err = h.cbStringList(nonEmptyList); err != nil {
					return err
				}
			}

		case ListOptionValueTypeUint32List:
			if v.Value == "" && h.ignoreZeroValue {
				continue
			}
			list := strings.Split(v.Value, ",")
			var intList []uint32
			for _, item := range list {
				x, err := strconv.ParseInt(item, 10, 32)
				if err != nil {
					return lberr.NewInvalidArg("invalid option value with type %d, expected uint32[]", v.Type)
				}
				intList = append(intList, uint32(x))
			}
			if h.cbUint32List != nil {
				if err = h.cbUint32List(intList); err != nil {
					return err
				}
			}

		case ListOptionValueTypeUint64:
			if v.Value == "" {
				continue
			}
			x, err := strconv.ParseUint(v.Value, 10, 64)
			if err != nil {
				return lberr.NewInvalidArg("invalid option value with type %d, expected uint64", v.Type)
			}
			if h.cbUint64 != nil {
				err = h.cbUint64(x)
				if err != nil {
					return err
				}
			}

		case ListOptionValueTypeTimeStampRange:
			var tStr []string
			if strings.Index(v.Value, ",") > 0 {
				tStr = strings.Split(v.Value, ",")
				if len(tStr) != 2 && h.ignoreZeroValue {
					continue
				}
			} else {
				tStr = strings.Split(v.Value, "-")
				if len(tStr) != 2 && h.ignoreZeroValue {
					continue
				}
			}
			t1, err := strconv.ParseUint(tStr[0], 10, 64)
			if err != nil {
				return lberr.NewInvalidArg("invalid option value with type %d, expected begin_time_stamp-end_time_stamp", v.Type)
			}
			t2, err := strconv.ParseUint(tStr[1], 10, 64)
			if err != nil {
				return lberr.NewInvalidArg("invalid option value with type %d, expected begin_time_stamp-end_time_stamp", v.Type)
			}
			if h.cbTimeStampRange != nil {
				if err = h.cbTimeStampRange(uint32(t1), uint32(t2)); err != nil {
					return err
				}
			}

		case ListOptionValueTypeUint64List:
			if v.Value == "" && h.ignoreZeroValue {
				continue
			}
			list := strings.Split(v.Value, ",")
			var intList []uint64
			for _, item := range list {
				x, err := strconv.ParseUint(item, 10, 64)
				if err != nil {
					return lberr.NewInvalidArg("invalid option value with type %d, expected uint64[]", v.Type)
				}
				intList = append(intList, x)
			}
			if h.cbUint64List != nil {
				if err = h.cbUint64List(intList); err != nil {
					return err
				}
			}

		case ListOptionValueTypeBool:
			//if v.Value != "0" && v.Value != "1" {
			//	continue
			//}
			value := strings.ToLower(v.Value)
			var x bool
			if utils.InSliceStr(value, []string{"1", "true"}) {
				x = true
			} else if utils.InSliceStr(value, []string{"0", "false"}) {
				x = false
			} else {
				continue
			}
			if h.cbBool != nil {
				err = h.cbBool(x)
				if err != nil {
					return err
				}
			}

		case ListOptionValueTypeInt32:
			if v.Value == "" {
				continue
			}
			x, err := strconv.ParseInt(v.Value, 10, 32)
			if err != nil {
				return lberr.NewInvalidArg("invalid option value with type %d, expected int32", v.Type)
			}
			if h.cbInt32 != nil {
				err = h.cbInt32(int32(x))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func toStr(i interface{}) string {
	t := reflect.TypeOf(i)
	k := t.Kind()
	switch k {
	case reflect.Int,
		reflect.Int32,
		reflect.Int64,
		reflect.Int16,
		reflect.Int8:
		return strconv.FormatInt(reflect.ValueOf(i).Int(), 10)
	case reflect.Uint,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint16,
		reflect.Uint8:
		return strconv.FormatUint(reflect.ValueOf(i).Uint(), 10)
	case reflect.String:
		return reflect.ValueOf(i).String()
	case reflect.Bool:
		if reflect.ValueOf(i).Bool() {
			return "1"
		} else {
			return "0"
		}
	}
	return fmt.Sprintf("%v", i)
}

func NewListOption(opts ...interface{}) *ListOption {
	if len(opts)%2 != 0 {
		log.Panicf("invalid number of opts argument %d", len(opts))
	}
	l := &ListOption{
		Limit:     defaultLimit,
		Offset:    0,
		SkipCount: false,
	}
	for i := 0; i < len(opts); i += 2 {
		l.AddOpt(opts[i], opts[i+1])
	}
	return l
}

func NewListOptionByPage(limit, offset string) *ListOption {
	newLimit, _ := strconv.ParseUint(limit, 10, 64)
	if newLimit == 0 {
		newLimit = defaultLimit
	}
	newOffset, _ := strconv.ParseUint(offset, 10, 64)
	return NewListOption().SetLimit(uint32(newLimit)).SetOffset(uint32(newOffset))
}

func (p *ListOption) SetLimit(limit uint32) *ListOption {
	p.Limit = limit
	return p
}

func (p *ListOption) SetOffset(offset uint32) *ListOption {
	p.Offset = offset
	return p
}

func (p *ListOption) IsOptExist(typ interface{}) bool {
	var typInt int32
	if reflect.TypeOf(typ).Kind() == reflect.Uint32 {
		typInt = int32(reflect.ValueOf(typ).Uint())
	} else {
		typInt = int32(reflect.ValueOf(typ).Int())
	}
	if typInt <= 0 {
		log.Panicf("invalid type %d", typ)
	}

	for _, opt := range p.Options {
		if opt.Type == typInt {
			return true
		}
	}

	return false
}

func (p *ListOption) GetOptValue(typ interface{}) (string, bool) {
	var typInt int32
	if reflect.TypeOf(typ).Kind() == reflect.Uint32 {
		typInt = int32(reflect.ValueOf(typ).Uint())
	} else {
		typInt = int32(reflect.ValueOf(typ).Int())
	}
	if typInt <= 0 {
		log.Panicf("invalid type %d", typ)
	}

	for _, opt := range p.Options {
		if opt.Type == typInt {
			return opt.Value, true
		}
	}

	return "", false
}

func (p *ListOption) AddOptIf(flag bool, typ, val interface{}) *ListOption {
	if flag {
		p.AddOpt(typ, val)
	}

	return p
}

func (p *ListOption) AddOpt(typ, val interface{}) *ListOption {
	var typInt int32
	if reflect.TypeOf(typ).Kind() == reflect.Uint32 {
		typInt = int32(reflect.ValueOf(typ).Uint())
	} else {
		typInt = int32(reflect.ValueOf(typ).Int())
	}
	if typInt <= 0 {
		log.Panicf("invalid type %d", typ)
	}
	typeOfVal := reflect.TypeOf(val)
	var strVal string
	if val == nil {
		strVal = ""
	} else {
		switch typeOfVal.Kind() {
		case reflect.Slice, reflect.Array:
			vv := reflect.ValueOf(val)
			n := vv.Len()
			var valList []string
			for j := 0; j < n; j++ {
				valList = append(valList, toStr(vv.Index(j).Interface()))
			}
			strVal = strings.Join(valList, ",")
		default:
			strVal = toStr(val)
		}
	}
	p.Options = append(p.Options,
		&ListOption_Option{Type: typInt, Value: strVal})
	return p
}

func (p *ListOption) SetSkipCount() *ListOption {
	p.SkipCount = true
	return p
}

func (p *ListOption) CloneSkipOpts() *ListOption {
	l := NewListOption().SetOffset(p.GetOffset()).SetLimit(p.GetLimit())
	if l.SkipCount {
		l.SetSkipCount()
	}
	return l
}

func getOptTypeFromInterface(typ interface{}) uint32 {
	t := reflect.TypeOf(typ)
	v := reflect.ValueOf(typ)
	if t.Kind() == reflect.Int32 {
		return uint32(v.Int())
	} else if t.Kind() == reflect.Uint32 {
		return uint32(v.Uint())
	} else {
		log.Panicf("unsupported type %s of opt with value %v", t.String(), typ)
	}
	return 0
}

func (p *ListOption) CloneChangeOptTypes(optPairs ...interface{}) *ListOption {
	l := p.CloneSkipOpts()
	if len(optPairs)%2 != 0 {
		log.Panicf("invalid number of opts argument %d", len(optPairs))
	}
	kv := map[uint32]uint32{}
	for i := 0; i < len(optPairs); i += 2 {
		typ := optPairs[i]
		val := optPairs[i+1]
		kv[getOptTypeFromInterface(typ)] = getOptTypeFromInterface(val)
	}
	for _, v := range p.Options {
		t := uint32(v.Type)
		if vv, ok := kv[t]; ok {
			l.AddOpt(vv, v.Value)
		}
	}
	return l
}
