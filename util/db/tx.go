package db

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

const (
	txKey int = iota
)

// GetTxFromCtx --
func GetTxFromCtx(c echo.Context) *gorm.DB {
	tx := c.Get(strconv.Itoa(txKey))
	if tx == nil {
		return db
	}
	return c.Get(strconv.Itoa(txKey)).(*gorm.DB)
}

// SetTxToCtx --
func SetTxToCtx(c echo.Context, tx *gorm.DB) {
	c.Set(strconv.Itoa(txKey), tx)
}
