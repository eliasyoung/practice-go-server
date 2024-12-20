package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	store *db.Store
}

func NewUserService(store *db.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) createUser(ctx context.Context, txq *db.Queries, username string, password []byte, email string) (db.CreateUserRow, error) {
	user := db.CreateUserParams{
		Username: username,
		Password: password,
		Email:    email,
	}

	row, err := txq.CreateUser(ctx, user)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return db.CreateUserRow{}, db.ErrDuplicateEmail
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)`:
			return db.CreateUserRow{}, db.ErrDuplicateUsername
		default:
			return db.CreateUserRow{}, err
		}
	}

	return row, nil
}

func (s *UserService) GetAllUser(ctx context.Context) ([]db.GetUsersRow, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	users, err := s.store.Queries.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, userId int64) (db.GetUserByIdRow, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	err := s.store.Queries.Follow(ctx, followParams)

	return err
}

func (s *UserService) UnfollowUserById(ctx context.Context, uid int64, ufid int64) error {
	unfollowParams := db.UnfollowParams{
		UserID:     uid,
		FollowerID: ufid,
	}

	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	err := s.store.Queries.Unfollow(ctx, unfollowParams)

	return err
}

func (s *UserService) CreateAndInviteUser(ctx context.Context, user db.CreateUserParams, token string, invitationExp time.Duration) error {
	return db.ExecWithTx(ctx, s.store, func(q *db.Queries) error {
		user, err := s.createUser(ctx, q, user.Username, user.Password, user.Email)
		if err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, q, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserService) createUserInvitation(ctx context.Context, q *db.Queries, token string, exp time.Duration, userID int64) error {
	params := db.CreateUserInvitationParams{
		Token:  []byte(token),
		UserID: userID,
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(exp),
			Valid: true,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	err := q.CreateUserInvitation(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Activate(ctx context.Context, token string) error {
	return db.ExecWithTx(ctx, s.store, func(q *db.Queries) error {
		user, err := s.getUserFromInvitation(ctx, q, token)
		if err != nil {
			return err
		}

		user.IsActive = true

		if err := s.updateUserInfoById(ctx, q, user.Username, user.Email, user.IsActive, user.ID); err != nil {
			return err
		}

		if err := s.deleteUserInvitationByUserId(ctx, q, user.ID); err != nil {
			return nil
		}

		return nil
	})
}

func (s *UserService) getUserFromInvitation(ctx context.Context, q *db.Queries, token string) (*db.GetUserFromInvitationRow, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ts := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}

	params := db.GetUserFromInvitationParams{
		Token:  []byte(hashToken),
		Expiry: ts,
	}

	user, err := q.GetUserFromInvitation(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, db.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserService) updateUserInfoById(ctx context.Context, q *db.Queries, username string, email string, is_active bool, id int64) error {
	param := db.UpdateUserInfoByIdParams{
		Username: username,
		Email:    email,
		IsActive: is_active,
		ID:       id,
	}

	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	if err := q.UpdateUserInfoById(ctx, param); err != nil {
		return err
	}

	return nil
}

func (s *UserService) deleteUserInvitationByUserId(ctx context.Context, q *db.Queries, uid int64) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	if err := q.DeleteUserInvitationByUserId(ctx, uid); err != nil {
		return err
	}

	return nil
}
