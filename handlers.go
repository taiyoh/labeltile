package labeltile

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

// HomeHandler provides controllers for toppage
type HomeHandler struct {
	container app.Container
}

// Top is toppage controller
func (h *HomeHandler) Top(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// OAuth2GoogleHandler provides controllers for oauth authorization by google
type OAuth2GoogleHandler struct {
	container app.Container
}

// Access is controller for oauth page redirection to google
func (h *OAuth2GoogleHandler) Access(c *gin.Context) {
	o := h.container.OAuth2Google()
	c.Redirect(http.StatusTemporaryRedirect, o.AuthCodeURL(""))
}

// Callback is controller for rediret page from google when authorized
func (h *OAuth2GoogleHandler) Callback(c *gin.Context) {}

// GraphQLHandler provides controllers for graphql requests
type GraphQLHandler struct {
	container app.Container
	schema    *graphql.Schema
}

// Run is controller for graphql request
func (h *GraphQLHandler) Run(c *gin.Context) {
	var err error
	var g *GraphQL
	req := c.Request
	if req.Header.Get("Content-Type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad Content-Type",
		})
		return
	}
	g, err = NewGraphQLRequest(h.schema, req.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := req.Context()
	if id, ok := c.Get("userID"); ok {
		ctx = context.WithValue(ctx, app.UserIDCtxKey, id)
	}
	res := g.Run(ctx)

	c.JSON(http.StatusOK, res)
}

// SetupHomeHandler provides initializing HomeHandler and binding its controller and endpoint
func SetupHomeHandler(router *gin.Engine, container app.Container) {
	h := &HomeHandler{container: container}
	router.GET("/", h.Top)
}

// SetupOAuth2GoogleHandler provides initializing OAuth2GoogleHandler and binding its controller and endpoint
func SetupOAuth2GoogleHandler(router *gin.Engine, container app.Container) {
	h := &OAuth2GoogleHandler{container: container}
	router.GET("/auth/google", h.Access)
	router.GET("/auth/google/callback", h.Callback)
}

// SetupGraphQLHandler provides initializing GraphQLHandler and binding its controller and endpoint
func SetupGraphQLHandler(router *gin.Engine, container app.Container) {
	h := GraphQLHandler{container: container, schema: InitializeGraphQLSchema(container)}
	router.POST("/graphql", h.Run)
}
