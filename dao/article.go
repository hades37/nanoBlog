package dao

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID          uint           `gorm:"primarykey"`                 // 主键ID
	Title       string         `gorm:"type:varchar(255);not null"` // 文章标题
	Slug        string         `gorm:"type:varchar(255);unique"`   // URL友好的标题别名
	Content     string         `gorm:"type:longtext;not null"`     // 文章内容
	Summary     string         `gorm:"type:varchar(500)"`          // 文章摘要
	CategoryID  uint           `gorm:"index"`                      // 分类ID
	Tags        string         `gorm:"type:varchar(255)"`          // 标签，以逗号分隔
	CoverImage  string         `gorm:"type:varchar(255)"`          // 封面图片
	AuthorID    uint           `gorm:"index;not null"`             // 作者ID
	Status      uint8          `gorm:"type:tinyint;default:1"`     // 状态 1:草稿 2:已发布 3:回收站
	IsTop       uint8          `gorm:"type:tinyint;default:0"`     // 是否置顶 0:否 1:是
	Views       uint           `gorm:"default:0"`                  // 浏览量
	CreatedAt   time.Time      `gorm:"type:datetime;not null"`     // 创建时间
	UpdatedAt   time.Time      `gorm:"type:datetime;not null"`     // 更新时间
	PublishedAt *time.Time     `gorm:"type:datetime"`              // 发布时间
	DeletedAt   gorm.DeletedAt `gorm:"index"`                      // 软删除
}

// TableName 设置表名
func (Article) TableName() string {
	return "article"
}

// CreateArticle 创建文章
func CreateArticle(db *gorm.DB, article *Article) error {
	return db.Create(article).Error
}

// GetArticleByID 根据ID获取文章
func GetArticleByID(db *gorm.DB, id uint) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticleBySlug 根据Slug获取文章
func GetArticleBySlug(db *gorm.DB, slug string) (*Article, error) {
	var article Article
	err := db.Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// UpdateArticle 更新文章
func UpdateArticle(db *gorm.DB, article *Article) error {
	return db.Save(article).Error
}

// DeleteArticle 删除文章（软删除）
func DeleteArticle(db *gorm.DB, id uint) error {
	return db.Delete(&Article{}, id).Error
}

// ListArticles 获取文章列表
func ListArticles(db *gorm.DB, page, pageSize int, conditions map[string]interface{}) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := db.Model(&Article{})

	// 添加查询条件
	for key, value := range conditions {
		query = query.Where(key, value)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = query.Order("is_top desc, created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// UpdateArticleStatus 更新文章状态
func UpdateArticleStatus(db *gorm.DB, id uint, status uint8) error {
	updates := map[string]interface{}{
		"status": status,
	}
	// 如果状态是已发布，则更新发布时间
	if status == 2 {
		updates["published_at"] = time.Now()
	}
	return db.Model(&Article{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementViews 增加文章浏览量
func IncrementViews(db *gorm.DB, id uint) error {
	return db.Model(&Article{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// GetArticlesByCategory 获取分类下的文章
func GetArticlesByCategory(db *gorm.DB, categoryID uint, page, pageSize int) ([]Article, int64, error) {
	return ListArticles(db, page, pageSize, map[string]interface{}{
		"category_id": categoryID,
		"status":      2, // 只获取已发布的文章
	})
}

// SearchArticles 搜索文章
func SearchArticles(db *gorm.DB, keyword string, page, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := db.Model(&Article{}).
		Where("status = ?", 2).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = query.Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
