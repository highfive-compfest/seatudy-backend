package course

import (
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Repository interface {
    GetAll(ctx context.Context) ([]Course, error)
    GetByID(ctx context.Context, id uuid.UUID) (Course, error)
    Create(ctx context.Context, course *Course) error
    Update(ctx context.Context, course *Course) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context) ([]Course, error) {
    var courses []Course
    if err := r.db.Preload("Materials.Attachments").Find(&courses).Error; err != nil {
        return nil, err
    }
    return courses, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (Course, error) {
    var course Course
    if err := r.db.Preload("Materials.Attachments").First(&course, "id = ?", id).Error; err != nil {
        return Course{}, err
    }
    return course, nil
}

func (r *repository) Create(ctx context.Context, course *Course) error {
    return r.db.WithContext(ctx).Create(course).Error
}

func (r *repository) Update(ctx context.Context, course *Course) error {
    return r.db.WithContext(ctx).Save(course).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&Course{}, id).Error
}
