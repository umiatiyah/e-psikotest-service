package query

func SqlQueryCek(tbl string) string {
	return `SELECT name, email, password FROM ` + tbl + ` WHERE email = $1`
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

func SqlGetCategoryIDFromQuestion() string {
	return `SELECT category_id FROM question WHERE id = $1`
}

func SqlGetQuestionIDFromAnswer() string {
	return `SELECT question_id FROM answer WHERE id = $1`
}

func SqlGetCurrentPassword(tbl string) string {
	return `SELECT password FROM ` + tbl + ` WHERE id = $1`
}

func SqlCreateTempTblBobot() string {
	return `drop table if exists tempBobot;
			CREATE TABLE tempBobot
			( userid int, category varchar (100), score int)
			;
			WITH BOBOT_CTE (userid, category, score)
			AS  
			(
			SELECT u.id, c.value, (SUM(a.score)) as bobot FROM history h JOIN category c ON h.category_id = c.id JOIN answer a ON h.answer_id = a.id JOIN users u ON h.user_id = u.id GROUP BY u.id, c.value ORDER BY bobot desc 
			)  
			insert into tempBobot
			SELECT *
			FROM BOBOT_CTE;`
}

func SqlGetMaxBobotCategory() string {
	return `select max(score) from tempBobot
			where category = $1`
}
