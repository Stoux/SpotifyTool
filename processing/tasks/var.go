package tasks

import (
	"context"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	ctx context.Context
)
