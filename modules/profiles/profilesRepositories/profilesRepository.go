package profilesrepositories

import (
	"fmt"

	"github.com/NattpkJsw/real-world-api-go/modules/profiles"
	"github.com/jmoiron/sqlx"
)

type IProfilesRepository interface {
	FindOneUserProfileByUsername(username string, curUserId int) (*profiles.Profile, error)
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
		fmt.Println("profile === ", profile)
		return nil, fmt.Errorf("user profile not found")
	}
	return profile, nil
}
