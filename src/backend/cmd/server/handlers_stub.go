package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/yourname/gym-management/internal/api" // ← go.mod の module 名に合わせる
)

type Server struct{}

// ----- Audit -----
func (s *Server) GetAuditLogs(c *gin.Context, params api.GetAuditLogsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Auth -----
func (s *Server) PostAuthLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostAuthRefresh(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ----- Dashboard -----
func (s *Server) GetDashboardNow(c *gin.Context, params api.GetDashboardNowParams) {
	c.JSON(http.StatusOK, gin.H{"groups": []gin.H{}})
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

// ----- ResourceGroups -----
func (s *Server) GetResourceGroups(c *gin.Context, params api.GetResourceGroupsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostResourceGroups(c *gin.Context, params api.PostResourceGroupsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PatchResourceGroupsGroupId(c *gin.Context, groupId string, params api.PatchResourceGroupsGroupIdParams) {
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

// ----- Sessions -----
func (s *Server) GetSessions(c *gin.Context, params api.GetSessionsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PostSessions(c *gin.Context, params api.PostSessionsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
func (s *Server) PatchSessionsSessionIdEnd(c *gin.Context, sessionId string, params api.PatchSessionsSessionIdEndParams) {
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
