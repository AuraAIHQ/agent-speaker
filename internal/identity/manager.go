// Package identity 提供身份管理功能
package identity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fiatjaf.com/nostr"
	"github.com/jason/agent-speaker/pkg/types"
)

const (
	configDirName  = ".agent-speaker"
	keystoreFile   = "keystore.json"
)

// Manager 管理身份
type Manager struct {
	configDir string
	keystore  *KeyStore
}

// KeyStore 是内部存储结构
type KeyStore struct {
	DefaultIdentity string                         `json:"default_identity"`
	Identities      map[string]*types.Identity     `json:"identities"`
	Secrets         map[string]string              `json:"secrets"` // nsec by nickname
	Contacts        map[string]*types.Contact      `json:"contacts"`
}

// NewManager 创建新的身份管理器
func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	configDir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}
	
	m := &Manager{
		configDir: configDir,
		keystore: &KeyStore{
			Identities: make(map[string]*types.Identity),
			Secrets:    make(map[string]string),
			Contacts:   make(map[string]*types.Contact),
		},
	}
	
	// 加载已有数据
	m.load()
	
	return m, nil
}

// Close 保存并关闭
func (m *Manager) Close() error {
	return m.save()
}

// Create 创建新身份
func (m *Manager) Create(nickname string) (*types.Identity, error) {
	if _, exists := m.keystore.Identities[nickname]; exists {
		return nil, fmt.Errorf("identity %s already exists", nickname)
	}
	
	// 生成密钥
	sk := nostr.Generate()
	pk := sk.Public()
	
	identity := &types.Identity{
		Nickname: nickname,
		Npub:     encodeNpub(pk),
		Created:  time.Now().Unix(),
	}
	
	m.keystore.Identities[nickname] = identity
	m.keystore.Secrets[nickname] = encodeNsec(sk)
	
	// 如果是第一个身份，设为默认
	if len(m.keystore.Identities) == 1 {
		m.keystore.DefaultIdentity = nickname
	}
	
	return identity, m.save()
}

// Get 获取身份
func (m *Manager) Get(nickname string) (*types.Identity, error) {
	if nickname == "" {
		nickname = m.keystore.DefaultIdentity
	}
	
	id, exists := m.keystore.Identities[nickname]
	if !exists {
		return nil, fmt.Errorf("identity %s not found", nickname)
	}
	
	return id, nil
}

// List 列出所有身份
func (m *Manager) List() ([]*types.Identity, error) {
	var result []*types.Identity
	for _, id := range m.keystore.Identities {
		result = append(result, id)
	}
	return result, nil
}

// SetDefault 设置默认身份
func (m *Manager) SetDefault(nickname string) error {
	if _, exists := m.keystore.Identities[nickname]; !exists {
		return fmt.Errorf("identity %s not found", nickname)
	}
	m.keystore.DefaultIdentity = nickname
	return m.save()
}

// 内部辅助函数

func (m *Manager) load() error {
	path := filepath.Join(m.configDir, keystoreFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, m.keystore)
}

func (m *Manager) save() error {
	path := filepath.Join(m.configDir, keystoreFile)
	data, err := json.MarshalIndent(m.keystore, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// encodeNpub 编码公钥为 npub
func encodeNpub(pk nostr.PubKey) string {
	// 简化版本，实际需要 bech32 编码
	return fmt.Sprintf("npub1%x", pk[:10])
}

// encodeNsec 编码私钥为 nsec  
func encodeNsec(sk nostr.SecretKey) string {
	// 简化版本
	return fmt.Sprintf("nsec1%x", sk[:10])
}
