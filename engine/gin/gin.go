package gin

import (
	"fmt"
	"net/http"

	ipfilter "github.com/Noah-Labs-Development/krakend-ipfilter"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

// Register register a ip filter middleware at gin
func Register(cfg *config.ServiceConfig, logger logging.Logger, engine *gin.Engine) {
	filterCfg := ipfilter.ParseConfig(cfg.ExtraConfig, logger)
	if filterCfg == nil {
		return
	}

	ipFilter := ipfilter.NewIPFilter(filterCfg)
	engine.Use(middleware(ipFilter, logger))
}

func middleware(ipFilter ipfilter.IPFilter, logger logging.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if ipFilter.Deny(ip) {
			logger.Error(fmt.Sprintf("krakend-ipfilter deny request from: %s", ip))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
