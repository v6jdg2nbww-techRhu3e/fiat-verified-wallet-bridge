package store

import (
    "sync"
    "time"
)

type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

type AuditLog struct {
    Event     string    `json:"event"`
    Actor     string    `json:"actor"`
    Detail    string    `json:"detail"`
    Timestamp time.Time `json:"timestamp"`
}

type Service struct {
    mu       sync.RWMutex
    users    map[string]*User
    auditLog []AuditLog
}

func NewService() *Service {
    return &Service{users: map[string]*User{}, auditLog: []AuditLog{}}
}

func (s *Service) UpsertUser(u *User) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.users[u.Email] = u
}

func (s *Service) LogAudit(event, actor, detail string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.auditLog = append(s.auditLog, AuditLog{
        Event:     event,
        Actor:     actor,
        Detail:    detail,
        Timestamp: time.Now().UTC(),
    })
}

func (s *Service) ListAudit() []AuditLog {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]AuditLog, len(s.auditLog))
    copy(out, s.auditLog)
    return out
}
