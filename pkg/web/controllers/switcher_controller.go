package controllers

import (
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type SwitcherController struct {
	config *config.SparrowConfig
}

func NewSwitcherController(config *config.SparrowConfig) *SwitcherController {
	return &SwitcherController{
		config: config,
	}
}
func (controller *SwitcherController) QueryNameSpace(c *gin.Context) {
	c.JSON(200, middleware.NewResponseSuccess(core.GetNs()))
}

func (controller *SwitcherController) QueryClass(c *gin.Context) {
	c.JSON(200, middleware.NewResponseSuccess(core.GetJSON(c.Query("ns"))))
}
