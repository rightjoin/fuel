package fuel

import "github.com/unrolled/render"

var VersionAfterPrefix = true

var rndr *render.Render

// HTTP Headers

var HeaderTotalRecords = "Fuel-Total-Records"
var HeaderPageSize = "Fuel-Page-Size"
var HeaderPageNum = "Fuel-Page-Num"
