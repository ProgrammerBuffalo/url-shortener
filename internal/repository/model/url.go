package model

import "github.com/google/uuid"

type UrlDao struct {
	Id    uuid.UUID
	Short string
	Long  string
}
