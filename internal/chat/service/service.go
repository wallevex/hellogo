package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.28.0"
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

	dsn = "root:root@tcp(192.168.3.3:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
)

type ChatServiceImpl struct {
	chatpb.UnimplementedChatServiceServer

	mu     *sync.RWMutex
	users  map[int64]*chatpb.User
	lastID int64

	db *sql.DB
}

func NewChatServiceImpl() (*ChatServiceImpl, error) {
	db, err := otelsql.Open("mysql", dsn, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		panic(err)
	}
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		panic(err)
	}
	s := &ChatServiceImpl{
		mu:     new(sync.RWMutex),
		users:  map[int64]*chatpb.User{},
		lastID: 101,
		db:     db,
	}

	return s, nil
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

func (s *ChatServiceImpl) SearchUsers(ctx context.Context, in *chatpb.SearchUsersRequest) (*chatpb.SearchUsersResponse, error) {
	tr := otel.Tracer("hello-chat")
	ctx, span := tr.Start(ctx, "SearchUsers")
	defer span.End()

	q := `SELECT * FROM user WHERE name LIKE ?`
	rows, err := s.db.QueryContext(ctx, q, in.Name+"%")
	if err != nil {
		return nil, err
	}

	var users []*chatpb.User
	for rows.Next() {
		var user chatpb.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return &chatpb.SearchUsersResponse{Users: users}, nil
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
