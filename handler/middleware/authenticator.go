package middleware

import (
	"github.com/kekim-go/Gateway/client"
	"github.com/kekim-go/Gateway/enum"
	grpc_author "github.com/kekim-go/Protobuf/gen/proto/author"
	"github.com/labstack/echo/v4"
)

func CheckAuth(authPool *client.Pool, c echo.Context) (*grpc_author.ApiAuthRes, enum.ResCode) {
	ctx := c.Request().Context()

	token, _ := c.Get("Token").(string)

	authConn, err := authPool.Get(ctx)
	defer authConn.Close()

	if err != nil {
		return nil, enum.InternalException
	}

	authorClient := grpc_author.NewApiAuthServiceClient(authConn)
	apiAuthRes, err := authorClient.Auth(ctx, &grpc_author.ApiAuthReq{
		NameSpace:    c.Param("nameSpace") + "/" + c.Param("version"),
		OperationUrl: c.Param("operation"),
		Token:        token,
	})

	var code enum.ResCode
	if err != nil {
		code = enum.Unknown
	} else {
		code = enum.FindResCode(apiAuthRes.Code)
	}

	return apiAuthRes, code
}
