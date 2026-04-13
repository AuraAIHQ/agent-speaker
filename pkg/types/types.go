// Package types 提供共享的数据类型定义
// 这些类型可以被外部项目导入使用
package types

import (
	"fiatjaf.com/nostr"
)

// Identity 表示用户身份（不包含私钥）
type Identity struct {
	Nickname string    `json:"nickname"`
	Npub     string    `json:"npub"`
	Created  int64     `json:"created"`
}

// Contact 表示联系人
type Contact struct {
	Nickname string    `json:"nickname"`
	Npub     string    `json:"npub"`
	AddedAt  int64     `json:"added_at"`
}

// Message 表示消息
type Message struct {
	ID            string    `json:"id"`
	SenderNpub    string    `json:"sender_npub"`
	RecipientNpub string    `json:"recipient_npub"`
	Content       string    `json:"content"`
	Plaintext     string    `json:"plaintext,omitempty"`
	CreatedAt     int64     `json:"created_at"`
	ReceivedAt    int64     `json:"received_at"`
	IsEncrypted   bool      `json:"is_encrypted"`
	IsIncoming    bool      `json:"is_incoming"`
}

// Encrypter 定义加密插件接口
type Encrypter interface {
	Name() string
	Encrypt(plaintext string, senderSK nostr.SecretKey, recipientPK nostr.PubKey) (string, error)
	Decrypt(ciphertext string, recipientSK nostr.SecretKey, senderPK nostr.PubKey) (string, error)
}

// Storage 定义存储插件接口
type Storage interface {
	SaveIdentity(identity *Identity) error
	GetIdentity(nickname string) (*Identity, error)
	ListIdentities() ([]*Identity, error)
	
	SaveContact(contact *Contact) error
	GetContact(nickname string) (*Contact, error)
	ListContacts() ([]*Contact, error)
	
	SaveMessage(msg *Message) error
	GetMessages(npub string, limit int) ([]*Message, error)
	GetConversation(npub1, npub2 string, limit int) ([]*Message, error)
}
