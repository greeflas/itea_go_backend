package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/greeflas/itea_go_backend/internal/model"
	"github.com/uptrace/bun"
)

type UserBunRepository struct {
	db *bun.DB
}

func NewUserBunRepository(db *bun.DB) *UserBunRepository {
	return &UserBunRepository{db: db}
}

func (r *UserBunRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	if err := r.db.NewSelect().Model(&users).Scan(ctx); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserBunRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)

	return err
}

func (r *UserBunRepository) GetById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := new(model.User)

	if err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserBunRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(ctx)

	return err
}

func (r *UserBunRepository) Delete(ctx context.Context, user *model.User) error {
	_, err := r.db.NewDelete().Model(user).WherePK().Exec(ctx)

	return err
}
