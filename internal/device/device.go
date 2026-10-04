package device

import "sync"

type Device struct {
    ID      string `json:"device_id"`
    Label   string `json:"label"`
    Email   string `json:"email"`
    Trusted bool   `json:"trusted"`
}

type Service struct {
    mu      sync.RWMutex
    devices map[string]*Device
}

func NewService() *Service {
    return &Service{devices: map[string]*Device{}}
}

func (s *Service) RegisterDevice(deviceID, label, email string, trusted bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.devices[deviceID] = &Device{ID: deviceID, Label: label, Email: email, Trusted: trusted}
}

func (s *Service) IsTrusted(deviceID, email string) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    d, ok := s.devices[deviceID]
    if !ok {
        return false
    }
    return d.Email == email && d.Trusted
}

func (s *Service) ListTrustedDevices() []Device {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]Device, 0)
    for _, d := range s.devices {
        if d.Trusted {
            out = append(out, *d)
        }
    }
    return out
}
