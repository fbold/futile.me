package util

import "github.com/sqids/sqids-go"

func GetId() *string {
	s, _ := sqids.New()
	id, _ := s.Encode([]uint64{1, 2, 3})
	return &id
}
