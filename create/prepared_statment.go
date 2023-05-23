package create

import "database/sql"

// CreatePreparedStatement creates an SQL statement for inserting
// records into the voters table.
func CreatePreparedStatement(tx *sql.Tx) (*sql.Stmt, error) {
	sqlString := CreateInsertSQL(selectedCols)
	stmt, err := tx.Prepare(sqlString)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
