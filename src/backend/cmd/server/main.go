package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"

	api "github.com/yourname/gym-management/internal/api" // ← module名に合わせて
)

// 開発用: bearerAuth を常に通す（最低限の形）
func authFunc(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	if input == nil {
		return nil
	}
	if input.SecuritySchemeName == "bearerAuth" {
		ah := input.RequestValidationInput.Request.Header.Get("Authorization")
		if ah == "" || !strings.HasPrefix(ah, "Bearer ") {
			return fmt.Errorf("unauthorized")
		}
		// TODO: 本実装ではここでJWT検証等を行う
		return nil
	}
	return nil
}

func main() {
	r := gin.Default()

	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("../openapi/openapi.yaml")
	if err != nil {
		log.Fatal(err)
	}

	r.Use(middleware.OapiRequestValidatorWithOptions(doc, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: authFunc,
		},
	}))

	s := &Server{}
	api.RegisterHandlersWithOptions(r, s, api.GinServerOptions{
	BaseURL: "/v1",
})

	log.Println("listening on :8080")
	_ = r.Run(":8080")
}
