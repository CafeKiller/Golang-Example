package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 定义Article结构体
type Article struct {
	// 继承Model结构体
	Model

	TagID      int    `json:"tag_id" gorm:"index"` // 定义TagID字段，类型为int，并设置为索引
	Tag        Tag    `json:"tag"`                 // 定义Tag字段，类型为Tag
	Title      string `json:"title"`               // 定义Title字段，类型为string
	Desc       string `json:"desc"`                // 定义Desc字段，类型为string
	Content    string `json:"content"`             // 定义Content字段，类型为string
	CreatedBy  string `json:"created_by"`          // 定义CreatedBy字段，类型为string
	ModifiedBy string `json:"modified_by"`         // 定义ModifiedBy字段，类型为string
	State      int    `json:"state"`               // 定义State字段，类型为int
}

// ExistArticleByID 通过ID判断文章是否存在
func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(article)
	if article.ID > 0 {
		return true
	}
	return false
}

// GetArticleTotal 获取文章总数
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

// GetArticles 获取多个文章
func GetArticles(pageNum int, pageSize int, maps interface{}) (article []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&article)
	return
}

// GetArticle 通过ID获取文章
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

// EditArticle 通过ID修改文章
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

// AddArticle 添加一篇文章
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})
	return true
}

// DeleteArticle 通过ID删除文章
func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}

// BeforeCreate 在创建之前调用
func (article *Article) BeforeCreate(scope gorm.Scope) error {
	// 设置创建时间
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

// BeforeUpdate 在更新之前调用
func (article *Article) BeforeUpdate(scope gorm.Scope) error {
	// 设置更新时间
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
