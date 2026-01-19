package util

import (
	"net/http"

	"github.com/a-h/templ"
)

// import "github.com/sqids/sqids-go"
//
// func GetId(source int) *string {
// 	s, _ := sqids.New()
// 	id, _ := s.Encode([]uint64{table, 2, 3})
// 	return &id
// }

func NullString(raw string) *string {
	var ptr *string
	if raw != "" {
		ptr = &raw
	}

	return ptr
}

func Serve(page func() templ.Component) func(w http.ResponseWriter, r *http.Request) {
	return templ.Handler(page()).ServeHTTP
}
