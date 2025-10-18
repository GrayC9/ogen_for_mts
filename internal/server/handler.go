//go:build ogen

package server

import (
	"context"
	"sync"

	api "github.com/example/ogen_for_mts/internal/api_1"
)

type inmem struct {
	mu    sync.RWMutex
	next  int64
	users map[int64]*api.User
}

func NewInMemoryHandler() (api.Handler, error) {
	return &inmem{users: make(map[int64]*api.User), next: 1}, nil
}

func (s *inmem) CreateUser(ctx context.Context, req *api.UserCreate) (api.CreateUserRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next
	s.next++
	u := &api.User{ID: id, Name: req.Name}
	if v, ok := req.Description.Get(); ok {
		u.Description.SetTo(v)
	}
	s.users[id] = u
	return u, nil
}

func (s *inmem) ListUsers(ctx context.Context) (api.ListUsersRes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]api.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, *u)
	}
	res := api.ListUsersOKApplicationJSON(out)
	return &res, nil
}

func (s *inmem) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[params.ID]
	if !ok {
		return &api.GetUserNotFound{}, nil
	}
	return u, nil
}

func (s *inmem) UpdateUser(ctx context.Context, req *api.UserUpdate, params api.UpdateUserParams) (api.UpdateUserRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[params.ID]
	if !ok {
		return &api.UpdateUserNotFound{}, nil
	}
	if req.Name.IsSet() {
		u.Name = req.Name.Or(u.Name)
	}
	if req.Description.IsSet() {
		if req.Description.IsNull() {
			u.Description.SetToNull()
		} else if v, ok := req.Description.Get(); ok {
			u.Description.SetTo(v)
		}
	}
	return u, nil
}

func (s *inmem) DeleteUser(ctx context.Context, params api.DeleteUserParams) (api.DeleteUserRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, params.ID)
	return &api.DeleteUserNoContent{}, nil
}
