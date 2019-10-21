package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EchoReq struct {
	Message string `form:"message"`
}

type EchoRes struct {
	Message string `form:"message" json:"message"`
}

func (s *Service) Echo(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &EchoReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	return req, &EchoRes{
		Message: req.Message,
	}, http.StatusOK, nil
}
