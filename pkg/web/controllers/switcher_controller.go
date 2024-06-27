package controllers

import (
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type SwitcherController struct {
	config  *config.SparrowConfig
	manager *core.SwitcherManager
}

func NewSwitcherController(config *config.SparrowConfig, manager *core.SwitcherManager) *SwitcherController {
	return &SwitcherController{
		config:  config,
		manager: manager,
	}
}
func (controller *SwitcherController) QueryNameSpace(c *gin.Context) {
	c.JSON(200, middleware.NewResponseSuccess(controller.manager.GetNs()))
}

func (controller *SwitcherController) QueryClass(c *gin.Context) {
	c.JSON(200, middleware.NewResponseSuccess(controller.manager.GetJSON(c.Query("ns"))))
}
