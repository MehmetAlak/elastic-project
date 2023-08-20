package rest

import (
	"elastic-project/interface/rest/docs"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

//go:generate swag init --parseDependency --parseInternal --parseDepth 1 -g server.go
// @title Elastic Search API
// @version 1.0
// @description An API that elastic search operations

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @Schemes http https

type server struct {
	elasticsearchEndpoint ElasticsearchEndpoint
}

type Server interface {
	SetupRouter() *gin.Engine
}

func NewServer(
	elasticsearchEndpoint ElasticsearchEndpoint) Server {
	return &server{
		elasticsearchEndpoint: elasticsearchEndpoint,
	}
}

func (server *server) SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.BestCompression))
	gin.SetMode(gin.ReleaseMode)

	//setUpNewRelic(router, server.newRelicConfig)

	if server.elasticsearchEndpoint != nil {
		router.PUT("/users/:id", server.elasticsearchEndpoint.Update())
		router.POST("/users", server.elasticsearchEndpoint.Create())
		router.GET("/users", server.elasticsearchEndpoint.Find())
		router.GET("/users-by", server.elasticsearchEndpoint.FindByKeyAndValue())
		router.GET("/users-by-query", server.elasticsearchEndpoint.FindByJsonQuery())
		router.DELETE("/users/:id", server.elasticsearchEndpoint.Delete())
	}

	//if server.healthEndpoint != nil {
	//	router.GET("/_monitoring/health", server.healthEndpoint.GetHealth())
	//}

	docs.SwaggerInfo.BasePath = "/"

	router.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("", func(c *gin.Context) {
		c.Redirect(302, "/swagger-ui/index.html")
		c.Abort()
	})
	router.GET("/swagger-ui.html", func(c *gin.Context) {
		c.Redirect(302, "/swagger-ui/index.html")
		c.Abort()
	})

	return router
}

//func setUpNewRelic(router *gin.Engine, newRelicConfig model.NewRelicConfig) {
//	cfg := newrelic.NewConfig(newRelicConfig.Name, "708418cc33cf343dc29fc9fbad5ca3756942d52a")
//	cfg.Logger = nrzap.Transform(logger.DefaultLogger())
//	cfg.Enabled = newRelicConfig.Enabled
//	cfg.DistributedTracer.Enabled = true
//	cfg.HighSecurity = false
//	cfg.SpanEvents.Enabled = true
//	cfg.TransactionEvents.Enabled = true
//
//	newRelicApp, err := newrelic.NewApplication(cfg)
//
//	if err != nil {
//		logger.DefaultLogger().Error("failed to make newRelicApp: " + err.Error())
//	}
//
//	if err = newRelicApp.WaitForConnection(time.Duration(30) * time.Second); nil != err {
//		logger.DefaultLogger().Error("failed to connect newRelic: " + err.Error())
//	} else {
//		newRelicMiddleware := nrgin.Middleware(newRelicApp)
//		router.Use(newRelicMiddleware)
//	}
//}
