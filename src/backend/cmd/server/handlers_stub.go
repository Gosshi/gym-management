// src/backend/cmd/server/handlers_stub.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/yourname/gym-management/internal/api" // ← go.mod の module 名に合わせて
)

// Server はここで一元定義（memStore を保持）
type Server struct {
	st *memStore
}

// ======== ここは未実装のまま（必要になったら実装する） ========

// ----- Auth -----
func (s *Server) PostAuthLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostAuthRefresh(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Gyms -----
func (s *Server) GetGymsMe(c *gin.Context, params api.GetGymsMeParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Members -----
func (s *Server) GetMembers(c *gin.Context, params api.GetMembersParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostMembers(c *gin.Context, params api.PostMembersParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Resources -----
func (s *Server) GetResources(c *gin.Context, params api.GetResourcesParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostResources(c *gin.Context, params api.PostResourcesParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PatchResourcesResourceId(c *gin.Context, resourceId string, params api.PatchResourcesResourceIdParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Zones -----
func (s *Server) GetZones(c *gin.Context, params api.GetZonesParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostZones(c *gin.Context, params api.PostZonesParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PatchZonesZoneId(c *gin.Context, zoneId string, params api.PatchZonesZoneIdParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Audit -----
func (s *Server) GetAuditLogs(c *gin.Context, params api.GetAuditLogsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
