package query

func SqlQueryCek(tbl string) string {
	return `SELECT name, email FROM ` + tbl + ` WHERE email = $1`
}

func SqlGetID(tbl string) string {
	return `SELECT id FROM ` + tbl + ` WHERE email = $1`
}

func SqlCount(tbl string) string {
	return `SELECT COUNT(*) FROM ` + tbl + ``
}

func SqlGetMaterialID(tbl string) string {
	return `SELECT id FROM ` + tbl + ` WHERE id = $1`
}

func SqlCekMaterialInOtherRelation(id int, column, tbl string) string {
	return `SELECT ` + column + ` FROM ` + tbl + ` WHERE ` + column + ` = $1`
}
