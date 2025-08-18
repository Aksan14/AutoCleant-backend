package util

import "database/sql"

func CommitOrRollBack(tx *sql.Tx) {
	err := recover() // tangkap panic jika ada
	if err != nil {
		_ = tx.Rollback()
		panic(err) // lempar ulang panic biar kelihatan errornya
	} else {
		_ = tx.Commit()
	}
}
