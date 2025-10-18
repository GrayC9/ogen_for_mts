//go:build ogen

package storage

import (
	"context"
	"sync"

	api "github.com/example/ogen_for_mts/internal/api_1"
)

type UserStorage interface {
	Create(ctx context.Context, req *api.UserCreate) (*api.User, error)
	List(ctx context.Context) ([]api.User, error)
	Get(ctx context.Context, id int64) (*api.User, bool, error)
	Update(ctx context.Context, id int64, req *api.UserUpdate) (*api.User, bool, error)
	Delete(ctx context.Context, id int64) error
}

type InMemory struct {
	mu    sync.RWMutex
	next  int64
	users map[int64]*api.User
}

func NewInMemory() *InMemory {
	return &InMemory{users: make(map[int64]*api.User), next: 1}
}

func (s *InMemory) Create(_ context.Context, req *api.UserCreate) (*api.User, error) {
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

func (s *InMemory) List(_ context.Context) ([]api.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]api.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, *u)
	}
	return out, nil
}

func (s *InMemory) Get(_ context.Context, id int64) (*api.User, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, false, nil
	}
	return u, true, nil
}

func (s *InMemory) Update(_ context.Context, id int64, req *api.UserUpdate) (*api.User, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	if !ok {
		return nil, false, nil
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
	return u, true, nil
}

func (s *InMemory) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, id)
	return nil
}
