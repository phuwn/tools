package middleware

import (
	"fmt"

	"github.com/labstack/echo"

	"github.com/phuwn/tools/db"
	"github.com/phuwn/tools/log"
)

// AddTransaction - middleware that help add transaction to handler
func AddTransaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != "GET" {
			tx := db.Get().Begin()
			db.SetTxToCtx(c, tx)
			err := next(c)

			msg, doRollBack := func() (string, bool) {
				if err != nil {
					return err.Error(), true
				}
				if rec := recover(); rec != nil {
					switch rec := rec.(type) {
					case error:
						return rec.Error(), true
					default:
						return fmt.Sprintf("unknown, recover: %v", rec), true
					}
				}
				if c.Response().Status == 500 {
					return "failed request", true
				}
				return "", false
			}()
			if doRollBack {
				log.Info("rollback transaction: ", msg)
				tx.Rollback()
				return err
			}
			tx.Commit()
			return err
		}
		return next(c)
	}
}
