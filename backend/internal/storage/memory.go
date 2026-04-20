package storage

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/amir-sheikh/secscan/backend/internal/model"
)

var ErrNotFound = errors.New("scan not found")

type MemoryStore struct {
	mu          sync.RWMutex
	scans       map[string]*model.Scan
	subscribers map[string]map[chan model.Scan]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		scans:       make(map[string]*model.Scan),
		subscribers: make(map[string]map[chan model.Scan]struct{}),
	}
}

func (s *MemoryStore) Create(scan *model.Scan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.scans[scan.ID] = cloneScan(scan)
	s.broadcastLocked(scan.ID, scan)
	return nil
}

func (s *MemoryStore) Get(id string) (*model.Scan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	scan, ok := s.scans[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneScan(scan), nil
}

func (s *MemoryStore) Update(id string, mutate func(*model.Scan) error) (*model.Scan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scan, ok := s.scans[id]
	if !ok {
		return nil, ErrNotFound
	}
	if err := mutate(scan); err != nil {
		return nil, err
	}

	updated := cloneScan(scan)
	s.scans[id] = updated
	s.broadcastLocked(id, updated)
	return cloneScan(updated), nil
}

func (s *MemoryStore) Subscribe(id string) (<-chan model.Scan, func(), error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scan, ok := s.scans[id]
	if !ok {
		return nil, nil, ErrNotFound
	}

	if _, exists := s.subscribers[id]; !exists {
		s.subscribers[id] = make(map[chan model.Scan]struct{})
	}

	ch := make(chan model.Scan, 8)
	s.subscribers[id][ch] = struct{}{}
	ch <- *cloneScan(scan)

	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if group, exists := s.subscribers[id]; exists {
			delete(group, ch)
			close(ch)
			if len(group) == 0 {
				delete(s.subscribers, id)
			}
		}
	}

	return ch, unsubscribe, nil
}

func (s *MemoryStore) broadcastLocked(id string, scan *model.Scan) {
	group, ok := s.subscribers[id]
	if !ok {
		return
	}

	for subscriber := range group {
		select {
		case subscriber <- *cloneScan(scan):
		default:
		}
	}
}

func cloneScan(scan *model.Scan) *model.Scan {
	if scan == nil {
		return nil
	}

	payload, _ := json.Marshal(scan)
	var out model.Scan
	_ = json.Unmarshal(payload, &out)
	return &out
}
