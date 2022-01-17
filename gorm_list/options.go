package gorm_list

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// QueryOptions ...
type QueryOptions func(o *queryOptions)

type queryOptions struct {
	IsFuzzy     bool
	SortBy      string // 默认升序，'-'开头则降序
	PageNum     int64
	PageSize    int64
	QueryMap    map[string]interface{}
	QueryString string
	QueryArgs   []interface{}
}

func (q *queryOptions) ExecuteOptions(opt []QueryOptions) {
	for _, fn := range opt {
		fn(q)
	}
}

const (
	defaultMaxPageSize = 10000
)

// GetDefaultQueryOptions ...
func GetDefaultQueryOptions() *queryOptions {
	return &queryOptions{
		IsFuzzy:  false,
		SortBy:   "",
		PageNum:  0,
		PageSize: 0,
	}
}

// WithPageNum 设置分页
func WithPageNum(pageNum int64) QueryOptions {
	return func(o *queryOptions) {
		o.PageNum = pageNum
	}
}

// WithPageSize 设置分页
func WithPageSize(pageSize int64) QueryOptions {
	return func(o *queryOptions) {
		o.PageSize = pageSize
	}
}

// WithSortBy 设置排序
func WithSortBy(sortBy string) QueryOptions {
	return func(o *queryOptions) {
		o.SortBy = sortBy
	}
}

// WithFuzzy 设置排序
func WithFuzzy(isFuzzy bool) QueryOptions {
	return func(o *queryOptions) {
		o.IsFuzzy = isFuzzy
	}
}

// WithQueryMap 设置查询序列
func WithQueryMap(queryMap map[string]interface{}) QueryOptions {
	return func(o *queryOptions) {
		o.QueryMap = queryMap
	}
}

// WithQueryString 设置自定义查询语句序列
func WithQueryString(queryString string, args ...interface{}) QueryOptions {
	return func(o *queryOptions) {
		o.QueryString = queryString
		o.QueryArgs = append(o.QueryArgs, args...)
	}
}

func MakeFuzzyQuery(queryMap map[string]interface{}) string {
	result := make([]string, 0, 0)
	for k, v := range queryMap {
		result = append(result, fmt.Sprintf("%s LIKE '%%%s%%'", k, v))
	}

	return strings.Join(result, " or ")
}

func (q *queryOptions) Limit() int {
	if q.PageSize == 0 {
		q.PageSize = defaultMaxPageSize
	}

	return int(q.PageSize)
}

// Offset .
func (q *queryOptions) Offset() int {
	if q.PageNum == 0 {
		q.PageNum = 1
	}
	if q.PageSize == 0 {
		q.PageSize = defaultMaxPageSize
	}
	return int((q.PageNum - 1) * q.PageSize)
}

type ListHandler struct {
	IsFuzzy  bool   `json:"is_fuzzy" form:"is_fuzzy"`
	SortBy   string `json:"sort_by" form:"sort_by"` // 默认升序，'-'开头则降序
	PageNum  int64  `json:"page_num" form:"page_num"`
	PageSize int64  `json:"page_size" form:"page_size"`
	QueryMap map[string]interface{}
}

const (
	maxOptionCount = 4
)

func (l *ListHandler) MakeQueryOptions() []QueryOptions {
	result := make([]QueryOptions, 0, maxOptionCount)
	result = append(result,
		WithPageSize(l.PageSize),
		WithFuzzy(l.IsFuzzy),
		WithPageNum(l.PageNum),
		WithSortBy(l.SortBy),
		WithQueryMap(l.QueryMap))

	return result
}

// List 通用List请求
func List(ctx context.Context, tr *gorm.DB, data interface{}, tableName string, options ...QueryOptions) (int64, error) {
	opts := GetDefaultQueryOptions()
	for _, option := range options {
		option(opts)
	}

	q := &gorm.DB{}
	if len(opts.QueryString) > 0 { // 自定义的query条件为最高优先级
		q = tr.Table(tableName).Where(opts.QueryString, opts.QueryArgs...)
	} else if opts.IsFuzzy {
		q = tr.Table(tableName).Where(MakeFuzzyQuery(opts.QueryMap))
	} else {
		q = tr.Table(tableName).Where(opts.QueryMap)
	}

	var total int64
	err := q.Count(&total).Error
	if err != nil {
		return 0, err
	}

	err = q.Order(opts.SortBy).Limit(opts.Limit()).Offset(opts.Offset()).Find(&data).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}
