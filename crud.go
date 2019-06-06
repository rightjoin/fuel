package fuel

import (
	"encoding/json"
	"fmt"

	"github.com/rightjoin/dorm"

	"github.com/jinzhu/gorm"
	"github.com/rightjoin/rutl/conv"
	"github.com/rightjoin/rutl/refl"
)

// Note: Even if there is a table column named PageSize, its
// SQL reference would be page_size (snake case). It would not
// conflict with default query values

var QPageSize = "page-size"
var QPageNum = "page-num"
var QOrderBy = "order-by"
var QOrderDir = "order-dir"
var QAltDb = "alt-db"

// FindHelper runs when model.Find GET service is invoked
func FindHelper(modl interface{}, ptrArrModel interface{}, ad Aide, dbo *gorm.DB) error {

	// If dbo is null, obtain
	// a DBO object
	if dbo == nil {
		dbo = QueryDB(ad)
	}

	flds := refl.NestedFields(modl)
	query := map[string]interface{}{}
	params := ad.Query()

	// If param matches one of model fields,
	// then only we use it in exact query
	for _, f := range flds {
		fldName := conv.CaseSnake(f.Name) // field's sql or db name

		if val, ok := params[fldName]; ok {
			query[fldName] = val
		}
	}

	// Total Record Count
	var count int
	err := dbo.Model(modl).Where(query).Count(&count).Error
	if err != nil {
		return err
	}
	ad.Response.Header().Set(HeaderTotalRecords, fmt.Sprintf("%d", count))

	// Pagination Size
	sizeVal, ok := params[QPageSize]
	size := conv.IntOr(sizeVal, 100)
	if size == -1 { // if pagination size is -1, then retreive all records
		size = count
	}
	ad.Response.Header().Set(HeaderPageSize, fmt.Sprintf("%d", size))

	// Page Number (to retreive)
	pageVal, ok := params[QPageNum]
	if !ok {
		pageVal = "1"
	}
	page := conv.IntOr(pageVal, 1)
	ad.Response.Header().Set(HeaderPageNum, fmt.Sprintf("%d", page))

	// Calculate Offset
	offset := (page - 1) * size

	// Order-By and Order Direction
	idDefault := false
	order, ok := params[QOrderBy]
	if !ok {
		// default order is "id"
		idDefault = true
		order = "id"
	}
	dirn, ok := params[QOrderDir]
	if !ok {
		if idDefault {
			dirn = "desc"
		} else {
			dirn = "asc"
		}
	}
	orderDir := fmt.Sprintf("%s %s", order, dirn)

	return dbo.Where(query).Order(orderDir).Offset(offset).Limit(size).Find(ptrArrModel).Error
}

// QueryHelper runs when model.Query POST service is invoked
func QueryHelper(modl interface{}, ptrArrModel interface{}, ad Aide, dbo *gorm.DB) error {

	// If dbo is null, obtain
	// a DBO object
	if dbo == nil {
		dbo = QueryDB(ad)
	}

	where := ad.Post()["where"]
	params := []interface{}{}
	if paramsStr, ok := ad.Post()["params"]; ok {
		err := json.Unmarshal([]byte(paramsStr), &params)
		if err != nil {
			return err
		}
	}

	// Record Count
	var count int
	err := dbo.Model(modl).Where(where, params).Count(&count).Error
	if err != nil {
		return err
	}
	ad.Response.Header().Set(HeaderTotalRecords, fmt.Sprintf("%d", count))

	return dbo.Where(where, params).Find(ptrArrModel).Error
}

// QueryDB sets up how fule gets the underlying ORM
// for the call. Default is to use master. However,
// if "alt-db" is present in Query String, then use
// a slave
func QueryDB(ad Aide) *gorm.DB {
	_, slave := ad.Query()[QAltDb]
	return dorm.GetORM(!slave)
}
