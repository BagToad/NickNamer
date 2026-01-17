package nicknames

import (
	"os"
	"testing"
)

func TestRememberAndForget(t *testing.T) {
	tmpFile := "/tmp/test_nicknames.json"
	defer os.Remove(tmpFile)

	m, err := NewManager(tmpFile)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	guildID := "test-guild-123"

	if m.HasWords(guildID) {
		t.Error("expected no words initially")
	}

	if !m.Remember("hello", guildID) {
		t.Error("expected to remember 'hello'")
	}

	if !m.HasWords(guildID) {
		t.Error("expected to have words after remember")
	}

	if m.Remember("hello", guildID) {
		t.Error("expected duplicate to fail")
	}

	words := m.GetWords(guildID)
	if len(words) != 1 || words[0] != "hello" {
		t.Errorf("expected [hello], got %v", words)
	}

	m.Forget("hello", guildID)
	if m.HasWords(guildID) {
		t.Error("expected no words after forget")
	}
}

func TestForgetAll(t *testing.T) {
	tmpFile := "/tmp/test_nicknames_all.json"
	defer os.Remove(tmpFile)

	m, err := NewManager(tmpFile)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	guildID := "test-guild-456"

	m.Remember("word1", guildID)
	m.Remember("word2", guildID)
	m.Remember("word3", guildID)

	m.ForgetAll(guildID)

	if m.HasWords(guildID) {
		t.Error("expected no words after forgetall")
	}
}

func TestGenerateName(t *testing.T) {
	tmpFile := "/tmp/test_nicknames_gen.json"
	defer os.Remove(tmpFile)

	m, err := NewManager(tmpFile)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	guildID := "test-guild-789"

	m.Remember("alpha", guildID)
	m.Remember("beta", guildID)
	m.Remember("gamma", guildID)

	name := m.GenerateName(guildID, 2)
	if name == "" {
		t.Error("expected non-empty name")
	}

	if m.GenerateName(guildID, 5) != "" {
		t.Error("expected empty name when requesting more words than available")
	}

	if m.GenerateName(guildID, 0) != "" {
		t.Error("expected empty name when requesting 0 words")
	}
}

func TestRoleName(t *testing.T) {
	tmpFile := "/tmp/test_nicknames_role.json"
	defer os.Remove(tmpFile)

	m, err := NewManager(tmpFile)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	guildID := "test-guild-role"

	if m.GetRoleName(guildID) != "NameChanger" {
		t.Error("expected default role name 'NameChanger'")
	}

	m.SetRoleName("CustomRole", guildID)
	if m.GetRoleName(guildID) != "CustomRole" {
		t.Error("expected role name 'CustomRole'")
	}
}

func TestPersistence(t *testing.T) {
	tmpFile := "/tmp/test_nicknames_persist.json"
	defer os.Remove(tmpFile)

	guildID := "test-guild-persist"

	m1, err := NewManager(tmpFile)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	m1.Remember("persistent", guildID)
	m1.SetRoleName("PersistRole", guildID)

	m2, err := NewManager(tmpFile)
	if err != nil {
		t.Fatalf("failed to create second manager: %v", err)
	}

	words := m2.GetWords(guildID)
	if len(words) != 1 || words[0] != "persistent" {
		t.Errorf("expected [persistent], got %v", words)
	}

	if m2.GetRoleName(guildID) != "PersistRole" {
		t.Error("expected role name 'PersistRole'")
	}
}
