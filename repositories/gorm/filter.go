package gormrepo

import (
	"fmt"
	"strings"

	"sisyphos/lib/peg"

	"gorm.io/gorm"
)

type Filter struct {
	db    *gorm.DB
	table string
}

func NewFilter(db *gorm.DB, table string) *Filter {
	return &Filter{db, table}
}

func (f *Filter) Filter(query string) ([]interface{}, error) {
	q, err := peg.Parse("", []byte(query))
	if err != nil {
		return nil, err
	}
	where := q.(string)

	dbq := fmt.Sprintf("Select `%s`.id from %s", f.table, f.table)
	if strings.Contains(where, "`tag`") {
		// TODO oh, please fix

		where = strings.ReplaceAll(where, "`tag`", "tags.text")
		dbq += " join actions_tags on actions.id = actions_tags.action_id join tags on actions_tags.tag_id = tags.id "
	}

	dbq = fmt.Sprintf("%s where %s", dbq, where)
	dbSess := f.db.Session(&gorm.Session{NewDB: true})

	idsStr := []string{}
	err = dbSess.Raw(dbq).Scan(&idsStr).Error
	if err != nil {
		return nil, err
	}
	ids := []interface{}{}
	for _, id := range idsStr {
		ids = append(ids, id)
	}

	return ids, nil
}
