package model

import "gorm.io/gorm"

// 聚合器Aggregator对应的数据库字段信息
type Email struct {
	Id      uint32 `gorm:"primaryKey;column:id" json:"id"`
	Address string `gorm:"column:address" json:"address"`
}

// 创建邮件地址至数据库
func (e Email) Create(db *gorm.DB) error {
	return db.Create(&e).Error
}

// 删除数据库中特定邮件地址
func (e Email) Delete(db *gorm.DB) error {
	return db.Where("id = ?", e.Id).Delete(&e).Error
}

// 获取数据库中特定邮件地址
func (e Email) Get(db *gorm.DB) (Email, error) {
	var Email Email
	err := db.Where("id = ?", e.Id).First(&Email).Error
	return Email, err

}

// 更新数据库中特定邮件地址
func (e Email) Update(db *gorm.DB) error {
	var temp Email
	temp.Id = e.Id
	return db.Model(&temp).Updates(e).Error
}

// 获取数据库中所有邮件地址
func (e Email) List(db *gorm.DB) ([]Email, error) {
	var Emails []Email
	err := db.Find(&Emails).Error
	if err != nil {
		return nil, err
	}
	return Emails, nil
}
