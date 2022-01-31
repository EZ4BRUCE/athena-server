package model

import "gorm.io/gorm"

// 在athena的service层要有个判断是否满足的方法
type Function struct {
	Id          uint32  `gorm:"primaryKey;column:id" json:"id"`
	Type        string  `gorm:"column:type" json:"type"`
	Threshold   float64 `gorm:"column:threshold" json:"threshold"`
	Description string  `gorm:"description" json:"description"`
}

func (f Function) Create(db *gorm.DB) error {
	return db.Create(&f).Error
}

func (f Function) Delete(db *gorm.DB) error {
	return db.Where("id = ?", f.Id).Delete(&f).Error
}

func (f Function) Get(db *gorm.DB) (Function, error) {
	var Function Function
	err := db.Where("id = ?", f.Id).First(&Function).Error
	return Function, err

}

func (f Function) Update(db *gorm.DB) error {
	var temp Function
	temp.Id = f.Id
	return db.Model(&temp).Updates(f).Error
}

func (f Function) List(db *gorm.DB) ([]Function, error) {
	var functions []Function
	err := db.Find(&functions).Error
	if err != nil {
		return nil, err
	}
	return functions, nil
}
