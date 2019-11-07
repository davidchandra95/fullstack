package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255";not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() []error {
	var err []error
	if p.Title == "" {
		err = append(err, errors.New("required title"))
	}

	if p.Content == "" {
		err = append(err, errors.New("required content"))
	}

	if p.AuthorID < 1 {
		err = append(err, errors.New("required author"))
	}

	return err
}

func (p *Post) Save(db *gorm.DB) (*Post, error) {
	var (
		err error
	)

	if err = db.Debug().Model(&Post{}).Create(&p).Error; err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		if err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error; err != nil {
			return &Post{}, err
		}
	}

	return p, err
}

func (p *Post) FindAll(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	if err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error; err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for _, post := range posts {
			if err := db.Debug().Model(&User{}).Where("id = ?", post.AuthorID).Take(&post.AuthorID).Error; err != nil {
				return &[]Post{}, err
			}
		}
	}

	return &posts, nil
}

func (p *Post) FindByID(db *gorm.DB, id uint64) (*Post, error) {
	var err error

	if err = db.Debug().Model(&Post{}).Where("id = ?", id).Take(&p).Error; err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		if err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error; err != nil {
			return &Post{}, err
		}
	}

	return p, err
}

func (p *Post) Update(db *gorm.DB, id uint64) (*Post, error) {
	var err error

	db = db.Debug().Model(&Post{}).Where("id = ?", id).Take(&Post{}).UpdateColumns(
		map[string]interface{}{
			"title": p.Title,
			"content": p.Content,
			"updated_at": time.Now(),
		},
	)
	if err = db.Debug().Model(&Post{}).Where("id = ?", id).Take(&p).Error; err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		if err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error; err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) Delete(db *gorm.DB, postID uint64, userID uint32) (int64, error) {
	if db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", postID, userID).Take(&Post{}).Delete(&Post{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("post not found")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}