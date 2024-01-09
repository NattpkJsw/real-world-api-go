package profiles

type Profile struct {
	Username  *string `db:"username"`
	Bio       *string `db:"bio"`
	Image     *string `db:"image"`
	Following *bool   `db:"following"`
}
