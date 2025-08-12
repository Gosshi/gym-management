// src/backend/cmd/server/handlers_impl.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/yourname/gym-management/internal/api" // ← module 名に合わせる
)

// Server は handlers_stub.go と同じ型を使う
// type Server struct { st *memStore } は既に定義済み想定

// -------- Dashboard --------
func (s *Server) GetDashboardNow(c *gin.Context, params api.GetDashboardNowParams) {
	gymID := params.XGymID

	// groups をスナップショット
	var groups []*api.ResourceGroup
	s.st.mu.RLock()
	for _, g := range s.st.getGroups(gymID) {
		groups = append(groups, g)
	}
	s.st.mu.RUnlock()

	items := make([]gin.H, 0, len(groups))
	for _, g := range groups {
		active := s.st.activeCount(gymID, g.Id)
		var wait *int
		if g.Capacity > 0 && active >= g.Capacity {
			w := 10 // TODO: 簡易推定。後で改善
			wait = &w
		}
		items = append(items, gin.H{
			"group_id":               g.Id,
			"name":                   g.Name,
			"capacity":               g.Capacity,
			"active_sessions":        active,
			"estimated_wait_minutes": wait,
		})
	}
	c.JSON(http.StatusOK, gin.H{"groups": items})
}

// -------- ResourceGroups --------
func (s *Server) GetResourceGroups(c *gin.Context, params api.GetResourceGroupsParams) {
	gymID := params.XGymID
	out := make([]api.ResourceGroup, 0)

	// スナップショット
	s.st.mu.RLock()
	for _, g := range s.st.getGroups(gymID) {
		out = append(out, *g)
	}
	s.st.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"items":       out,
		"next_cursor": nil,
	})
}

func (s *Server) PostResourceGroups(c *gin.Context, params api.PostResourceGroupsParams) {
	var body api.PostResourceGroupsJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	rules := map[string]interface{}{}
	if body.RulesJson != nil {
		rules = *body.RulesJson
	}
	g := s.st.createGroup(params.XGymID, body.Name, body.ZoneId, body.Capacity, rules)
	c.JSON(http.StatusCreated, g)
}

func (s *Server) PatchResourceGroupsGroupId(c *gin.Context, groupId string, params api.PatchResourceGroupsGroupIdParams) {
	var body api.PatchResourceGroupsGroupIdJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	s.st.mu.Lock()
	defer s.st.mu.Unlock()
	m := s.st.getGroups(params.XGymID)
	g := m[groupId]
	if g == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	if body.Name != nil {
		g.Name = *body.Name
	}
	if body.ZoneId != nil {
		g.ZoneId = body.ZoneId
	}
	if body.Capacity != nil {
		g.Capacity = *body.Capacity
	}
	if body.RulesJson != nil {
		if len(*body.RulesJson) == 0 {
			g.RulesJson = nil
		} else {
			tmp := *body.RulesJson
			g.RulesJson = &tmp
		}
	}
	c.JSON(http.StatusOK, g)
}

// -------- Sessions --------
func (s *Server) PostSessions(c *gin.Context, params api.PostSessionsParams) {
	var body api.PostSessionsJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	ses, err := s.st.startSession(params.XGymID, body.GroupId, body.MemberId, body.ResourceId, body.Source)
	if err != nil {
		if err.Error() == "capacity_exceeded" {
			c.JSON(http.StatusConflict, gin.H{"error": "capacity_exceeded"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ses)
}

func (s *Server) PatchSessionsSessionIdEnd(c *gin.Context, sessionId string, params api.PatchSessionsSessionIdEndParams) {
	var body api.PatchSessionsSessionIdEndJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		// ended_at は任意なので、ボディが無い場合はOK
		body = api.PatchSessionsSessionIdEndJSONRequestBody{}
	}
	ses, err := s.st.endSession(params.XGymID, sessionId, body.EndedAt)
	if err != nil {
		if err.Error() == "not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ses)
}

func (s *Server) GetSessions(c *gin.Context, params api.GetSessionsParams) {
	gymID := params.XGymID

	// スナップショット
	list := make([]*api.UtilizationSession, 0)
	s.st.mu.RLock()
	for _, v := range s.st.getSessions(gymID) {
		list = append(list, v)
	}
	s.st.mu.RUnlock()

	from := params.From
	to := params.To
	gid := params.GroupId
	mid := params.MemberId

	out := make([]api.UtilizationSession, 0)
	for _, ses := range list {
		if gid != nil && ses.GroupId != *gid {
			continue
		}
		if mid != nil && (ses.MemberId == nil || *ses.MemberId != *mid) {
			continue
		}
		if from != nil && ses.StartedAt.Before((*from).UTC()) {
			continue
		}
		if to != nil && ses.StartedAt.After((*to).UTC()) {
			continue
		}
		out = append(out, *ses)
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       out,
		"next_cursor": nil,
	})
}

// -------- ここより下は未実装のままでもOK（スタブが返る） --------

// 例: Members/Resources/Zones などは既存の 501 スタブが動作
// 必要になったらここに実装を足していく
