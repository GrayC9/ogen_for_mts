//go:build ogen

package server

import (
	"context"

	api "github.com/example/ogen_for_mts/internal/api_1"
	"github.com/example/ogen_for_mts/internal/storage"
)

type handler struct {
	store storage.UserStorage
}

func NewInMemoryHandler() (api.Handler, error) {
	return &handler{store: storage.NewInMemory()}, nil
}

func (h *handler) CreateUser(ctx context.Context, req *api.UserCreate) (api.CreateUserRes, error) {
	u, err := h.store.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (h *handler) ListUsers(ctx context.Context) (api.ListUsersRes, error) {
	list, err := h.store.List(ctx)
	if err != nil {
		return nil, err
	}
	res := api.ListUsersOKApplicationJSON(list)
	return &res, nil
}

func (h *handler) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	u, ok, err := h.store.Get(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &api.GetUserNotFound{}, nil
	}
	return u, nil
}

func (h *handler) UpdateUser(ctx context.Context, req *api.UserUpdate, params api.UpdateUserParams) (api.UpdateUserRes, error) {
	u, ok, err := h.store.Update(ctx, params.ID, req)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &api.UpdateUserNotFound{}, nil
	}
	return u, nil
}

func (h *handler) DeleteUser(ctx context.Context, params api.DeleteUserParams) (api.DeleteUserRes, error) {
	if err := h.store.Delete(ctx, params.ID); err != nil {
		return nil, err
	}
	return &api.DeleteUserNoContent{}, nil
}
