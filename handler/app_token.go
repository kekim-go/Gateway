package handler

import (
	"github.com/kekim-go/Gateway/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegistAppToken(c echo.Context) error {
	c.Request().Context()
	appToken := &models.AppToken{}
	if err := c.Bind(appToken); err != nil {
		return err
	}

	return nil
}
