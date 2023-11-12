package models

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

//// BeforeCreate GORM V2中使用BeforeCreate和BeforeUpdate钩子的方式有所变化
//func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {
//	tag.CreatedOn = time.Now().Unix()
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(tx *gorm.DB) (err error) {
//	tag.ModifiedOn = time.Now().Unix()
//	return nil
//}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	return tag.ID > 0
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	return tag.ID > 0
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetTagTotal(maps interface{}) (count int64) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Updates(data)
	return true
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})
	return true
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}

func CleanAllTag() bool {
	result := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})
	return result.Error == nil
}
