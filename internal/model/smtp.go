package model

import "gorm.io/gorm"

type Smtp struct {
	Id       uint32 `gorm:"primaryKey;column:id" json:"id"`
	Host     string `gorm:"column:host" json:"host"`
	Port     int    `gorm:"column:port" json:"port"`
	IsSSL    bool   `gorm:"column:isSSL" json:"isSSL"`
	UserName string `gorm:"column:user_name" json:"userName"`
	Password string `gorm:"column:password" json:"password"`
	From     string `gorm:"column:from" json:"from"`
}

// 创建告警规则至数据库
func (s Smtp) Create(db *gorm.DB) error {
	return db.Create(&s).Error
}

// 删除数据库中特定告警规则
func (s Smtp) Delete(db *gorm.DB) error {
	return db.Where("id = ?", s.Id).Delete(&s).Error
}

// 获取数据库中特定告警规则
func (s Smtp) Get(db *gorm.DB) (Smtp, error) {
	var Smtp Smtp
	err := db.Where("id = ?", s.Id).First(&Smtp).Error
	return Smtp, err

}

// 更新数据库中特定告警规则
func (s Smtp) Update(db *gorm.DB) error {
	var temp Smtp
	temp.Id = s.Id
	return db.Model(&temp).Updates(s).Error
}

// 获取数据库中所有告警规则
func (s Smtp) List(db *gorm.DB) ([]Smtp, error) {
	var Smtps []Smtp
	err := db.Find(&Smtps).Error
	if err != nil {
		return nil, err
	}
	return Smtps, nil
}
