package helperfunctions

import "database/sql"

func CommitOrRollBack(tx *sql.Tx) {
	err := recover()

	if err != nil {
		errRollback := tx.Rollback()
		if err != nil {
			panic(errRollback)
		}

		panic(err)
	} else {
		errCommit := tx.Commit()
		if err != nil {
			panic(errCommit)
		}
	}
}
