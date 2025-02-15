package handler

import (
	"fmt"
	"net/http"

	grpc_author "github.com/kekim-go/Protobuf/gen/proto/author"
	"github.com/kekim-go/Gateway/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RegisterPortalKey(c echo.Context) error {
	logrus.WithFields(logrus.Fields{
		"action": "RegisterPortalKey",
	}).Info("call poratl key")

	token, _ := c.Get("Token").(string)
	logrus.WithFields(logrus.Fields{
		"token": token,
	}).Info("call poratl key")

	if token != "data-go-kr-prod-token" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"msg": "허가되지 않은 접근 입니다.",
		})
	}

	appToken := new(models.AppToken)
	if err := c.Bind(appToken); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"msg": "알 수 없는 오류가 발생하였습니다. 문제가 지속되면 관리자에게 문의하세요",
		})
	}

	ctx := c.Request().Context()

	authConn, err := h.authPool.Get(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Exception in get grpc connection")

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"msg": "Author 연계 과정 중 알 수 없는 오류가 발생하였습니다. 문제가 지속되면 관리자에게 문의하세요",
		})
	}
	defer authConn.Close()
	appTokenClient := grpc_author.NewAppTokenManagerClient(authConn)
	appTokenRes, err := appTokenClient.Create(ctx, &grpc_author.AppTokenReq{
		NameSpace: appToken.NameSpace,
		Token:     appToken.Token,
	})

	if err != nil || appTokenRes.Code != grpc_author.AppTokenRes_OK {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Exception in portal token regist")

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"msg": fmt.Sprintf("%s, %s", appTokenRes.Message, err.Error()),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"msg": "정상처리 되었습니다.",
	})
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
