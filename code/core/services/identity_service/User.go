package identity_service

import (
	"context"
	"errors"

	"github.com/spacetimi/timi_shared_server/v2/code/config"
	"github.com/spacetimi/timi_shared_server/v2/code/core/services/storage_service"
	"github.com/spacetimi/timi_shared_server/v2/code/core/services/storage_service/storage_typedefs"
)

const kUBVersion = 1

// Implements IBlob
type UserBlob struct {
	UserId        int64
	CreatedTime   int64
	LastLoginTime int64

	UserName                 string
	UserEmailAddress         string
	UserEmailAddressVerified bool

	storage_typedefs.BlobDescriptor `bson:"ignore"`
}

func newUserBlob(userId int64) *UserBlob {
	user := UserBlob{
		UserId: userId,
	}
	user.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_SHARED,
		config.GetAppName()+"::user",
		[]string{"UserId"},
		kVersion,
		true)
	return &user
}

func loadUserBlobByUserId(userId int64, ctx context.Context) (*UserBlob, error) {
	user := newUserBlob(userId)

	err := storage_service.GetBlobByPrimaryKeys(user, ctx)
	if err != nil {
		return nil, errors.New("error getting user blob: " + err.Error())
	}

	return user, nil
}
