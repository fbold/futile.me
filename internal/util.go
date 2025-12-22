package util

import "github.com/sqids/sqids-go"

func GetId(source int) *string {
	s, _ := sqids.New()
	id, _ := s.Encode([]uint64{table, 2, 3})
	return &id
}
