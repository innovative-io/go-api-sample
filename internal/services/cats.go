package services

import (
	"context"
	"errors"

	"github.com/innovative-io/go-api-sample/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CatsService interface {
	Add(ctx context.Context, cat *models.Cat) (*uuid.UUID, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context) ([]models.Cat, error)
	GetOne(ctx context.Context, id uuid.UUID) (*models.Cat, error)
	Update(ctx context.Context, id uuid.UUID, cat *models.Cat) error
}

type catsService struct {
	db *gorm.DB
}

func NewCatsService(db *gorm.DB) CatsService {
	return &catsService{
		db: db,
	}
}

func (s *catsService) Add(ctx context.Context, cat *models.Cat) (*uuid.UUID, error) {
	if err := s.db.Create(cat).Error; err != nil {
		return nil, err
	}
	return &cat.ID, nil
}

func (s *catsService) Delete(ctx context.Context, id uuid.UUID) error {
	db := s.db.Delete(&models.Cat{}, id)
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected < 1 {
		return ErrNotFound
	}
	return nil
}

func (s *catsService) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Cat{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *catsService) Get(ctx context.Context) ([]models.Cat, error) {
	cats := make([]models.Cat, 0)
	if err := s.db.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *catsService) GetOne(ctx context.Context, id uuid.UUID) (*models.Cat, error) {
	cat := new(models.Cat)
	if err := s.db.First(cat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return cat, nil
}

func (s *catsService) Update(ctx context.Context, id uuid.UUID, cat *models.Cat) error {
	db := s.db.Model(&models.Cat{ID: id}).Updates(models.Cat{
		Name:      cat.Name,
		Breed:     cat.Breed,
		Color:     cat.Color,
		Birthdate: cat.Birthdate,
		Weight:    cat.Weight,
	})
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected < 1 {
		return ErrNotFound
	}
	return nil
}
