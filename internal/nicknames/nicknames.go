package nicknames

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

// GuildData holds the words and role configuration for a single guild.
type GuildData struct {
	Words []string `json:"words"`
	Role  string   `json:"NICKNAMER_ROLE"`
}

// Manager handles nickname word storage and generation.
type Manager struct {
	mu       sync.RWMutex
	dataFile string
	data     map[string]*GuildData
}

// NewManager creates a new nickname manager with the specified data file.
func NewManager(dataFile string) (*Manager, error) {
	m := &Manager{
		dataFile: dataFile,
		data:     make(map[string]*GuildData),
	}
	if err := m.Load(); err != nil {
		m.data = make(map[string]*GuildData)
		if err := m.Save(); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// ensureGuild initializes guild data if it doesn't exist.
func (m *Manager) ensureGuild(guildID string) {
	if _, exists := m.data[guildID]; !exists {
		m.data[guildID] = &GuildData{
			Words: []string{},
			Role:  "NameChanger",
		}
	}
}

// Remember adds a word to the guild's word pool.
func (m *Manager) Remember(word, guildID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureGuild(guildID)
	wordLower := strings.ToLower(word)
	for _, w := range m.data[guildID].Words {
		if strings.ToLower(w) == wordLower {
			return false
		}
	}
	m.data[guildID].Words = append(m.data[guildID].Words, word)
	if err := m.Save(); err != nil {
		log.Printf("failed to save data: %v", err)
	}
	return true
}

// Forget removes a word from the guild's word pool.
func (m *Manager) Forget(word, guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[guildID]; !exists {
		return
	}
	words := m.data[guildID].Words
	for i, w := range words {
		if w == word {
			m.data[guildID].Words = append(words[:i], words[i+1:]...)
			if err := m.Save(); err != nil {
				log.Printf("failed to save data: %v", err)
			}
			return
		}
	}
}

// ForgetAll clears all words for the guild.
func (m *Manager) ForgetAll(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[guildID]; !exists {
		return
	}
	m.data[guildID].Words = []string{}
	if err := m.Save(); err != nil {
		log.Printf("failed to save data: %v", err)
	}
}

// GetWords returns the list of words for a guild.
func (m *Manager) GetWords(guildID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if gd, exists := m.data[guildID]; exists {
		return gd.Words
	}
	return []string{}
}

// HasWords checks if the guild has any words stored.
func (m *Manager) HasWords(guildID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	gd, exists := m.data[guildID]
	if !exists {
		return false
	}
	return len(gd.Words) > 0
}

// GetRoleName returns the role name for a guild.
func (m *Manager) GetRoleName(guildID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if gd, exists := m.data[guildID]; exists && gd.Role != "" {
		return gd.Role
	}
	return "NameChanger"
}

// SetRoleName sets the role name for a guild.
func (m *Manager) SetRoleName(roleName, guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureGuild(guildID)
	m.data[guildID].Role = roleName
	if err := m.Save(); err != nil {
		log.Printf("failed to save data: %v", err)
	}
}

// GenerateName creates a random nickname using n words from the pool.
func (m *Manager) GenerateName(guildID string, n int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	gd, exists := m.data[guildID]
	if !exists || len(gd.Words) < n || n < 1 {
		return ""
	}

	used := make(map[int]bool)
	parts := make([]string, 0, n)

	for len(parts) < n {
		idx := rand.Intn(len(gd.Words))
		if !used[idx] {
			used[idx] = true
			parts = append(parts, gd.Words[idx])
		}
	}
	return strings.Join(parts, " ")
}

// Save writes the data to disk.
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.dataFile, data, 0644)
}

// Load reads the data from disk.
func (m *Manager) Load() error {
	if _, err := os.Stat(m.dataFile); os.IsNotExist(err) {
		return nil
	}
	data, err := os.ReadFile(m.dataFile)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &m.data)
}
