package model

import "gorm.io/gorm"

// 聚合函数和告警规则分来存，保证扩展性
// 蕴蓄同一个等级不同行为（描述不同），id区分
type Rule struct {
	Id          uint32 `gorm:"primaryKey;column:id" json:"id"`
	Level       string `gorm:"level" json:"level"`
	Action      string `gorm:"action" json:"action"`
	Description string `gorm:"description" json:"description"`
}

func (r Rule) Create(db *gorm.DB) error {
	return db.Create(&r).Error
}

func (r Rule) Delete(db *gorm.DB) error {
	return db.Where("id = ?", r.Id).Delete(&r).Error
}

func (r Rule) Get(db *gorm.DB) (Rule, error) {
	var Rule Rule
	err := db.Where("id = ?", r.Id).First(&Rule).Error
	return Rule, err

}

func (r Rule) Update(db *gorm.DB) error {
	var temp Rule
	temp.Id = r.Id
	return db.Model(&temp).Updates(r).Error
}

func (r Rule) List(db *gorm.DB) ([]Rule, error) {
	var rules []Rule
	err := db.Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}
