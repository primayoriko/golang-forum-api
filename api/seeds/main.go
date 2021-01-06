package seeds

import (
	"gitlab.com/hydra/forum-api/api"
)

func SeedData() error {
	db, err := api.ConnectDB()
	if err != nil {
		return err
	}

	SeedUsers(db)
	SeedThreads(db)
	SeedPosts(db)

	return err
}
