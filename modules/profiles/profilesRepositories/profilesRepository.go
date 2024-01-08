package profilesrepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/profiles"
	"github.com/jmoiron/sqlx"
)

type IProfilesRepository interface {
	FindOneUserProfileByUsername(username string, curUserId int) (*profiles.Profile, error)
	FollowUser(username string, curUserId int) (*profiles.Profile, error)
}

type profilesRepository struct {
	db *sqlx.DB
}

func ProfilesRepository(db *sqlx.DB) IProfilesRepository {
	return &profilesRepository{
		db: db,
	}
}

func (r *profilesRepository) FindOneUserProfileByUsername(username string, curUserId int) (*profiles.Profile, error) {
	query := `
		SELECT
		"username",
		"bio",
		"image",
		(
			SELECT
			(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
			FROM "user_follows" "uf"
			WHERE "uf"."follower_id" = $2 AND 
			"uf"."following_id" = (SELECT
									"id"
									FROM "users"
									WHERE "username" = $1)
		) AS "following"
		FROM "users"
		WHERE "username" = $1;`

	profile := new(profiles.Profile)
	if err := r.db.Get(profile, query, username, curUserId); err != nil {
		return nil, fmt.Errorf("user profile not found")
	}
	return profile, nil
}

func (r *profilesRepository) FollowUser(username string, curUserId int) (*profiles.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "user_follows"
	(
		"following_id",
		"follower_id"
	)
	SELECT
    (
		SELECT 
		"id" 
		FROM "users" 
		WHERE "username" = $1),$2
	WHERE 
	(
    	SELECT 
			"id" 
		FROM "users" 
		WHERE "username" = $1
	) != $2;
	`
	if _, err := r.db.ExecContext(ctx, query, username, curUserId); err != nil {
		return nil, fmt.Errorf("fail to follow the user")
	}

	return r.FindOneUserProfileByUsername(username, curUserId)
}
