package forum

import (
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type IRepository interface {
	CreateDiscussion(discussion *schema.ForumDiscussion) error
	GetDiscussionByID(id uuid.UUID) (*schema.ForumDiscussion, error)
	GetDiscussionsByCourseID(courseID uuid.UUID, page int, limit int) ([]*schema.ForumDiscussion, int64, error)
	UpdateDiscussion(discussion *schema.ForumDiscussion) error
	DeleteDiscussion(id uuid.UUID) error

	CreateReply(reply *schema.ForumReply) error
	GetReplyByID(id uuid.UUID) (*schema.ForumReply, error)
	GetRepliesByDiscussionID(discussionID uuid.UUID, page int, limit int) ([]*schema.ForumReply, int64, error)
	UpdateReply(reply *schema.ForumReply) error
	DeleteReply(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &repository{db: db}
}

func (r *repository) CreateDiscussion(discussion *schema.ForumDiscussion) error {
	return r.db.Create(discussion).Error
}

func (r *repository) GetDiscussionByID(id uuid.UUID) (*schema.ForumDiscussion, error) {
	var discussion schema.ForumDiscussion
	err := r.db.Where("id = ?", id).First(&discussion).Error
	return &discussion, err
}

func (r *repository) GetDiscussionsByCourseID(courseID uuid.UUID, page int, limit int) ([]*schema.ForumDiscussion, int64, error) {
	var discussions []*schema.ForumDiscussion
	var total int64

	tx := r.db.Model(&schema.ForumDiscussion{}).Where("course_id = ?", courseID)

	tx.Count(&total)

	tx.Order("id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&discussions)

	return discussions, total, nil
}

func (r *repository) UpdateDiscussion(discussion *schema.ForumDiscussion) error {
	tx := r.db.Updates(discussion)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) DeleteDiscussion(id uuid.UUID) error {
	tx := r.db.Where("id = ?", id).Delete(&schema.ForumDiscussion{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) CreateReply(reply *schema.ForumReply) error {
	return r.db.Create(reply).Error
}

func (r *repository) GetReplyByID(id uuid.UUID) (*schema.ForumReply, error) {
	var reply schema.ForumReply
	err := r.db.Where("id = ?", id).First(&reply).Error
	return &reply, err
}

func (r *repository) GetRepliesByDiscussionID(discussionID uuid.UUID, page int, limit int) ([]*schema.ForumReply, int64, error) {
	var replies []*schema.ForumReply
	var total int64

	tx := r.db.Model(&schema.ForumReply{}).Where("forum_discussion_id = ?", discussionID)

	tx.Count(&total)

	tx.Order("id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&replies)

	return replies, total, nil
}

func (r *repository) UpdateReply(reply *schema.ForumReply) error {
	tx := r.db.Updates(reply)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) DeleteReply(id uuid.UUID) error {
	tx := r.db.Where("id = ?", id).Delete(&schema.ForumReply{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
