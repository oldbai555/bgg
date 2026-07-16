package iam

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserThirdPartyRepository interface {
	FindByOpenID(ctx context.Context, provider, openID string) (*iammodel.AdminUserThirdParty, error)
	Create(ctx context.Context, bind *iammodel.AdminUserThirdParty) error
}

type userThirdPartyRepository struct {
	model iammodel.AdminUserThirdPartyModel
	conn  sqlx.SqlConn
}

func NewUserThirdPartyRepository(repo *repository.Repository) UserThirdPartyRepository {
	return &userThirdPartyRepository{model: repo.AdminUserThirdPartyModel, conn: repo.DB}
}

func (r *userThirdPartyRepository) FindByOpenID(ctx context.Context, provider, openID string) (*iammodel.AdminUserThirdParty, error) {
	bind, err := r.model.FindOneByProviderOpenId(ctx, provider, openID)
	if err == iammodel.ErrNotFound {
		return nil, nil
	}
	return bind, err
}

func (r *userThirdPartyRepository) Create(ctx context.Context, bind *iammodel.AdminUserThirdParty) error {
	result, err := r.model.Insert(ctx, bind)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	bind.Id = uint64(id)
	return nil
}
