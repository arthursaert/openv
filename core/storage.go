package core

import (
	"github.com/google/uuid"
)

func GenerateRepoID() string {
	return uuid.New().String()
}

func GenerateCommitID() string {
	return uuid.New().String()
}
