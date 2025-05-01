package postgres

import (
	"context"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/db"
)

type postgresTaskRepository struct {
	Conn db.DB
}

func NewPostgresTaskRepository(conn db.DB) domain.TaskRepository {
	return &postgresTaskRepository{conn}
}

func (p postgresTaskRepository) GetByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresTaskRepository) Store(ctx context.Context, task *domain.Task) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresTaskRepository) SetWordSelected(ctx context.Context, taskId domain.TaskID, dictId domain.DictionaryID) error {
	//TODO implement me
	panic("implement me")
}
