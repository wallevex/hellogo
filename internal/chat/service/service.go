package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	chatpb "hellogo/api/chat"
	"hellogo/internal/chat/service/algo"
	"sync"
)

const (
	defaultID   = 100
	defaultName = "default"
)

type ChatServiceImpl struct {
	chatpb.UnimplementedChatServiceServer

	mu     *sync.RWMutex
	users  map[int64]*chatpb.User
	lastID int64
}

func NewChatServiceImpl() *ChatServiceImpl {
	return &ChatServiceImpl{
		mu: new(sync.RWMutex),
		users: map[int64]*chatpb.User{
			100: {
				Id:   defaultID,
				Name: defaultName,
			},
		},
		lastID: 101,
	}
}

func (s *ChatServiceImpl) GetUser(ctx context.Context, in *chatpb.GetUserRequest) (*chatpb.User, error) {
	s.mu.RLock()
	user, ok := s.users[in.Id]
	s.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user(id=%v)", in.Id)
	}

	return user, nil
}
func (s *ChatServiceImpl) CreateUser(ctx context.Context, in *chatpb.CreateUserRequest) (*chatpb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	in.User.Id = s.lastID
	s.users[s.lastID] = in.User
	s.lastID++

	return in.User, nil
}

func (s *ChatServiceImpl) DeleteUser(ctx context.Context, in *chatpb.DeleteUserRequest) (*emptypb.Empty, error) {
	if in.Id == defaultID {
		return nil, status.Errorf(codes.PermissionDenied, "prohibited from deleting default user")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.users[in.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user[id=%d]", in.Id)
	}
	delete(s.users, in.Id)

	return &emptypb.Empty{}, nil
}

func (s *ChatServiceImpl) Ask(ctx context.Context, in *chatpb.AskRequest) (*chatpb.AskResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[in.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user[id=%d]", in.Id)
	}

	res := algo.Fibo(in.N)
	text := fmt.Sprintf("hello %s, fibo(%d) is %d.", user.Name, in.N, res)

	return &chatpb.AskResponse{Text: text}, nil
}
