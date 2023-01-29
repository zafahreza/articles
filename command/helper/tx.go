package helper

import "gorm.io/gorm"

func CommitOrRollback(tx *gorm.DB) {
	err := recover()
	if err != nil {
		rollback := tx.Rollback()
		if rollback.Error != nil {
			panic(err)
		}
		panic(err)
	} else {
		commit := tx.Commit()
		if commit.Error != nil {
			panic(err)
		}
	}
}
