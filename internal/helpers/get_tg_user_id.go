package helpers

import "gorm.io/gorm"

func GetTgUserID(db *gorm.DB, userID uint) (int64, error) {
	var res *int64

	err := db.Raw("SELECT tg_user_id FROM users WHERE id = ?", userID).Scan(&res).Error

	if err != nil {
		return 0, err
	}

	if res == nil {
		return 0, nil
	}

	return *res, nil
}
