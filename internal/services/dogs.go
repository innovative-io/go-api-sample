package services

import (
	"context"
	"errors"

	"github.com/innovative-io/go-api-sample/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DogsService interface {
	Add(ctx context.Context, dog *models.Dog) (*uuid.UUID, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context) ([]models.Dog, error)
	GetOne(ctx context.Context, id uuid.UUID) (*models.Dog, error)
	Update(ctx context.Context, id uuid.UUID, dog *models.Dog) error
}

type dogsService struct {
	db *gorm.DB
}

func NewDogsService(db *gorm.DB) DogsService {
	return &dogsService{
		db: db,
	}
}

func (s *dogsService) Add(ctx context.Context, dog *models.Dog) (*uuid.UUID, error) {
	if err := s.db.Create(dog).Error; err != nil {
		return nil, err
	}
	return &dog.ID, nil
}

func (s *dogsService) Delete(ctx context.Context, id uuid.UUID) error {
	db := s.db.Delete(&models.Dog{}, id)
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected < 1 {
		return ErrNotFound
	}
	return nil
}

func (s *dogsService) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Dog{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *dogsService) Get(ctx context.Context) ([]models.Dog, error) {
	dogs := make([]models.Dog, 0)
	if err := s.db.Find(&dogs).Error; err != nil {
		return nil, err
	}
	return dogs, nil
}

func (s *dogsService) GetOne(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
	dog := new(models.Dog)
	if err := s.db.First(dog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dog, nil
}

func (s *dogsService) Update(ctx context.Context, id uuid.UUID, dog *models.Dog) error {
	db := s.db.Model(&models.Dog{ID: id}).Updates(models.Dog{
		Name:      dog.Name,
		Breed:     dog.Breed,
		Color:     dog.Color,
		Birthdate: dog.Birthdate,
		Weight:    dog.Weight,
	})
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected < 1 {
		return ErrNotFound
	}
	return nil
}
