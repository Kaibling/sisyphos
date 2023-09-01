package pagination

import (
	"fmt"
	"strings"
)

type IDReader interface {
	ReadPaginationIds(sqlQuery string) ([]DBResult, error)
}

func New(sortInfo SortInfo, dbReader IDReader, dbInfo DatabaseInfo) ([]string, *Cursor, error) {
	//validate data
	if sortInfo.Table == "" {
		return nil, nil, fmt.Errorf("no table configured")
	}
	if sortInfo.Field == "" {
		tInfo, ok := dbInfo.Tables[sortInfo.Table]
		if !ok {
			return nil, nil, fmt.Errorf("table %s not configured", sortInfo.Table)
		}
		sortInfo.Field = tInfo.PrimaryField
	}
	if sortInfo.Order == "" {
		sortInfo.Order = "ASC"
	} else {
		if strings.EqualFold(sortInfo.Order, "ASC") {
			sortInfo.Order = "ASC"
		}
		if strings.EqualFold(sortInfo.Order, "DESC") {
			sortInfo.Order = "DESC"
		}
	}
	var ids []string
	var err error
	var c *Cursor
	if sortInfo.CursorInfo == nil {
		ids, c, err = noCursor(dbReader, sortInfo, dbInfo)
		if err != nil {
			return nil, nil, err
		}
	} else {
		ids, c, err = withCursor(dbReader, sortInfo, dbInfo)
		if err != nil {
			return nil, nil, err
		}
	}
	c.Finish()
	return ids, c, nil

}

func noCursor(dbReader IDReader, sortInfo SortInfo, dbInfo DatabaseInfo) ([]string, *Cursor, error) {
	selectedIds := sortInfo.Field
	orderBy := sortInfo.Field
	primaryId := dbInfo.Tables[sortInfo.Table].PrimaryField
	multiField := false
	if primaryId != sortInfo.Field {
		orderBy = fmt.Sprintf("%s, %s", sortInfo.Field, primaryId)
		selectedIds = fmt.Sprintf("%s, %s", sortInfo.Field, primaryId)
		multiField = true
	}
	q := fmt.Sprintf("SELECT %s FROM user ORDER BY %s %s limit %d", selectedIds, orderBy, sortInfo.Order, sortInfo.Limit+1)
	c := &Cursor{}
	fmt.Println(q)

	dbresult, err := dbReader.ReadPaginationIds(q)
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(dbresult) > sortInfo.Limit {
		dbresult = dbresult[:len(dbresult)-1]
		cids := []string{dbresult[sortInfo.Limit-1].PID}
		if multiField {
			cids = append(cids, ConvertToString(dbresult[sortInfo.Limit-1].SortID))
		}
		c.CreateAfter(sortInfo.Field, cids...)
	}
	primaryIDs := []string{}
	for _, dbRes := range dbresult {
		primaryIDs = append(primaryIDs, dbRes.PID)
	}
	return primaryIDs, c, nil
}

func withCursor(dbReader IDReader, sortInfo SortInfo, dbInfo DatabaseInfo) ([]string, *Cursor, error) {
	order := sortInfo.Order
	var directionOperator string
	var min = 0
	var max = sortInfo.Limit

	// TODO ????
	if sortInfo.Order == "ASC" {
		directionOperator = ">"
	} else {
		directionOperator = "<"
	}
	if sortInfo.CursorInfo.Direction == "BEFORE" {
		order = "DESC"
		if sortInfo.Order == "DESC" {
			order = "ASC"
			directionOperator = ">"
		} else {
			directionOperator = "<"
		}
	}
	selectedIds := sortInfo.Field + " as pid"
	orderBy := sortInfo.Field
	primaryId := dbInfo.Tables[sortInfo.Table].PrimaryField
	multiField := false
	outerSelectIds := "pid"
	innerWhere := sortInfo.Field + " " + directionOperator + " " + sortInfo.CursorInfo.PrimaryId
	if primaryId != sortInfo.Field {
		orderBy = fmt.Sprintf("%s, %s", sortInfo.Field, primaryId)
		selectedIds = fmt.Sprintf("%s as sid, %s as pid", sortInfo.Field, primaryId)
		multiField = true
		outerSelectIds += ",sid"
		innerWhere = sortInfo.CursorInfo.SortField + " " + directionOperator + " '" + sortInfo.CursorInfo.SortId + "'"
		innerWhere += " or ( " + sortInfo.CursorInfo.SortField + " =  '" + sortInfo.CursorInfo.SortId + "' AND " + primaryId + " " + directionOperator + " " + sortInfo.CursorInfo.PrimaryId + "   )"

	}
	q := "SELECT " + outerSelectIds +
		" FROM (select " + selectedIds +
		" from " + sortInfo.Table +
		" where " + innerWhere +
		" ORDER BY " + orderBy + " " + order +
		" limit " + fmt.Sprintf("%d", sortInfo.Limit+1) + "  )" +
		" ORDER BY " + outerSelectIds + " " + sortInfo.Order
	fmt.Println(q)
	c := &Cursor{}

	dbresult, err := dbReader.ReadPaginationIds(q)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(dbresult)
	if sortInfo.CursorInfo.Direction == "AFTER" {
		if len(dbresult) > sortInfo.Limit {
			dbresult = dbresult[:len(dbresult)-1]
			cids := []string{dbresult[sortInfo.Limit-1].PID}
			if multiField {
				cids = append(cids, ConvertToString(dbresult[sortInfo.Limit-1].SortID))
			}
			c.CreateAfter(sortInfo.Field, cids...)
		}
		cids := []string{dbresult[min].PID}
		if multiField {
			cids = append(cids, ConvertToString(dbresult[min].SortID))
		}
		c.CreateBefore(sortInfo.Field, cids...)
	} else {
		if len(dbresult) > sortInfo.Limit {
			dbresult = dbresult[1:]
			cids := []string{dbresult[min].PID}
			if multiField {
				cids = append(cids, ConvertToString(dbresult[min].SortID))
			}
			c.CreateBefore(sortInfo.Field, cids...)
		}
		cids := []string{dbresult[max-1].PID}
		if multiField {
			cids = append(cids, ConvertToString(dbresult[max-1].SortID))
		}
		c.CreateAfter(sortInfo.Field, cids...)
	}
	primaryIDs := []string{}
	for _, dbRes := range dbresult {
		primaryIDs = append(primaryIDs, dbRes.PID)
	}
	return primaryIDs, c, nil
}
