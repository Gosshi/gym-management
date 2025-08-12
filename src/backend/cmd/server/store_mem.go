// src/backend/cmd/server/store_mem.go
package main

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	api "github.com/yourname/gym-management/internal/api" // ← go.mod の module 名に合わせる
)

type memStore struct {
	mu       sync.RWMutex
	// gym_id -> group_id -> group
	groups   map[string]map[string]*api.ResourceGroup
	// gym_id -> session_id -> session
	sessions map[string]map[string]*api.UtilizationSession
}

func newStore() *memStore {
	return &memStore{
		groups:   make(map[string]map[string]*api.ResourceGroup),
		sessions: make(map[string]map[string]*api.UtilizationSession),
	}
}

func newID(prefix string) string {
	var b [6]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%s_%x", prefix, b[:])
}

func (s *memStore) getGroups(gymID string) map[string]*api.ResourceGroup {
	if s.groups[gymID] == nil {
		s.groups[gymID] = make(map[string]*api.ResourceGroup)
	}
	return s.groups[gymID]
}

func (s *memStore) getSessions(gymID string) map[string]*api.UtilizationSession {
	if s.sessions[gymID] == nil {
		s.sessions[gymID] = make(map[string]*api.UtilizationSession)
	}
	return s.sessions[gymID]
}

func (s *memStore) activeCount(gymID, groupID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.sessions[gymID]
	if m == nil {
		return 0
	}
	n := 0
	for _, ses := range m {
		if ses.GroupId == groupID && ses.Status == api.UtilizationSessionStatusActive {
			n++
		}
	}
	return n
}

func (s *memStore) capacityOf(gymID, groupID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g := s.groups[gymID]
	if g == nil || g[groupID] == nil {
		return 0
	}
	return g[groupID].Capacity
}

func (s *memStore) createGroup(gymID, name string, zoneID *string, capacity int, rules map[string]interface{}) *api.ResourceGroup {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := newID("grp")
	r := &api.ResourceGroup{
		Id:       id,
		GymId:    gymID,
		Name:     name,
		ZoneId:   zoneID,
		Capacity: capacity,
	}
	if len(rules) > 0 {
		r.RulesJson = &rules
	}
	if s.groups[gymID] == nil {
		s.groups[gymID] = make(map[string]*api.ResourceGroup)
	}
	s.groups[gymID][id] = r
	return r
}

func (s *memStore) startSession(gymID, groupID string, memberID, resourceID *string, source *api.PostSessionsJSONBodySource) (*api.UtilizationSession, error) {
	// capacity check
	capacity := s.capacityOf(gymID, groupID)
	if capacity > 0 && s.activeCount(gymID, groupID) >= capacity {
		return nil, fmt.Errorf("capacity_exceeded")
	}

	now := time.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()

	id := newID("ses")
	status := api.UtilizationSessionStatusActive
	var src api.UtilizationSessionSource
	if source != nil {
		switch *source {
		case api.PostSessionsJSONBodySourceApi:
			src = api.UtilizationSessionSourceApi
		case api.PostSessionsJSONBodySourceStaff:
			src = api.UtilizationSessionSourceStaff
		default:
			src = api.UtilizationSessionSourceSelf
		}
	} else {
		src = api.UtilizationSessionSourceSelf
	}
	ses := &api.UtilizationSession{
		Id:         id,
		GymId:      gymID,
		GroupId:    groupID,
		MemberId:   memberID,
		ResourceId: resourceID,
		StartedAt:  now,
		Status:     status,
		Source:     &src,
	}
	if s.sessions[gymID] == nil {
		s.sessions[gymID] = make(map[string]*api.UtilizationSession)
	}
	s.sessions[gymID][id] = ses
	return ses, nil
}

func (s *memStore) endSession(gymID, sessionID string, endedAt *time.Time) (*api.UtilizationSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.sessions[gymID]
	if m == nil || m[sessionID] == nil {
		return nil, fmt.Errorf("not_found")
	}
	ses := m[sessionID]
	if ses.Status != api.UtilizationSessionStatusActive {
		return ses, nil // idempotent
	}
	end := time.Now().UTC()
	if endedAt != nil {
		end = endedAt.UTC()
	}
	ses.EndedAt = &end
	dur := int(end.Sub(ses.StartedAt).Seconds())
	ses.DurationSec = &dur
	ses.Status = api.UtilizationSessionStatusEnded
	return ses, nil
}
