package storage

import (
	"context"
	"github.com/SerjLeo/storage_bot/internal/models"
)

type Storage interface {
	Save(ctx context.Context, p *models.Page) error
	Remove(ctx context.Context, p *models.Page) error
	Pick(ctx context.Context, username string) (*models.Page, error)
	IsExist(ctx context.Context, p *models.Page) (bool, error)
	List(ctx context.Context, username string) ([]*models.Page, error)
	MarkAsSeen(ctx context.Context, p *models.Page) error
	DeleteSeen(ctx context.Context, username string) error
}

var NoPagesFoundError error
