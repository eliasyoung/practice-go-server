package service

import (
	"context"
	"time"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
)

type UserService struct {
	store *db.Store
}

func NewUserService(store *db.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) CreateUser(ctx context.Context, txq *db.Queries, username string, password []byte, email string) (db.CreateUserRow, error) {
	user := db.CreateUserParams{
		Username: username,
		Password: password,
		Email:    email,
	}

	row, err := txq.CreateUser(ctx, user)
	if err != nil {
		switch {
		case err.Error() == `pg: duplicate key value violates unique constraint "user_email_key"`:
			return db.CreateUserRow{}, db.ErrDuplicateEmail
		case err.Error() == `pg: duplicate key value violates unique constraint "user_username_key"`:
			return db.CreateUserRow{}, db.ErrDuplicateUsername
		default:
			return db.CreateUserRow{}, err
		}
	}

	return row, nil
}

func (s *UserService) GetAllUser(ctx context.Context) ([]db.GetUsersRow, error) {

	users, err := s.store.Queries.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, userId int64) (db.GetUserByIdRow, error) {
	user, err := s.store.Queries.GetUserById(ctx, userId)
	if err != nil {
		return db.GetUserByIdRow{}, err
	}

	return user, nil
}

func (s *UserService) FollowUserById(ctx context.Context, uid int64, fid int64) error {
	followParams := db.FollowParams{
		UserID:     uid,
		FollowerID: fid,
	}

	err := s.store.Queries.Follow(ctx, followParams)

	return err
}

func (s *UserService) UnfollowUserById(ctx context.Context, uid int64, ufid int64) error {

	unfollowParams := db.UnfollowParams{
		UserID:     uid,
		FollowerID: ufid,
	}

	err := s.store.Queries.Unfollow(ctx, unfollowParams)

	return err
}

func (s *UserService) CreateAndInviteUser(ctx context.Context, user db.CreateUserParams, token string, invitationExp time.Duration) error {
	return db.ExecWithTx(ctx, s.store, func(q *db.Queries) error {
		_, err := s.CreateUser(ctx, q, user.Username, user.Password, user.Email)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserService) createUserInvitation(ctx context.Context, token string, exp time.Duration, userID int64) error {

	return nil
}
