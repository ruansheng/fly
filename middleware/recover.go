package middleware

import (
	"fmt"
	"github.com/ruansheng/fly"
)

func Recover() fly.MiddlewareFunc {
	return func(next fly.HandlerFunc) fly.HandlerFunc {
		return func(ctx fly.Context) error {
			fmt.Println("middleware->Recover", ctx.Path())
			defer func() {
				var err error
				if r := recover(); r != nil {
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					fmt.Println("middleware->recover info", err.Error())
				}
			}()
			return next(ctx)
		}
	}
}
