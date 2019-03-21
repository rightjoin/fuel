package fuel

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rightjoin/dorm"
	"github.com/rightjoin/rutl/conv"
	"github.com/rightjoin/rutl/refl"
)

// FindHelper runs when model.Find GET service is invoked
func FindHelper(modl interface{}, ptrArrModel interface{}, ad Aide, dbo *gorm.DB) error {

	// If dbo is null, obtain
	// access to read-only replica
	if dbo == nil {
		dbo = dorm.GetORM(false)
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

	// Count
	var count int
	err := dbo.Model(modl).Where(query).Count(&count).Error
	if err != nil {
		return err
	}
	ad.Response.Header().Set(dorm.HeaderCount, fmt.Sprintf("%d", count))

	// Size
	sizeVal, ok := params[":page-size"]
	if !ok {
		sizeVal = "-1"
	}
	size := conv.IntOr(sizeVal, 100)

	// Page Number
	pageVal, ok := params[":page-num"]
	if !ok {
		pageVal = "1"
	}
	page := conv.IntOr(pageVal, 1)

	// Calculate Offset
	offset := 0
	if size != -1 {
		offset = (page - 1) * size
	}

	// Order
	order, ok := params[":order"]
	if !ok {
		order = "id" // default order is "id"
	}

	return dbo.Where(query).Order(order).Offset(offset).Limit(size).Find(ptrArrModel).Error
}

// QueryHelper runs when model.Query POST service is invoked
func QueryHelper(modl interface{}, ptrArrModel interface{}, ad Aide, dbo *gorm.DB) error {

	// If dbo is null, obtain
	// access to read-only replica
	if dbo == nil {
		dbo = dorm.GetORM(false)
	}

	where := ad.Post()["where"]
	params := []interface{}{}
	if paramsStr, ok := ad.Post()["params"]; ok {
		err := json.Unmarshal([]byte(paramsStr), &params)
		if err != nil {
			return err
		}
	}

	// Count
	var count int
	err := dbo.Model(modl).Where(where, params).Count(&count).Error
	if err != nil {
		return err
	}
	ad.Response.Header().Set(dorm.HeaderCount, fmt.Sprintf("%d", count))

	return dbo.Where(where, params).Find(ptrArrModel).Error
}
