-- name: CreateUserInvitation :exec
INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3);

-- name: DeleteUserInvitationByUserId :exec
DELETE FROM user_invitations WHERE user_id = $1;
