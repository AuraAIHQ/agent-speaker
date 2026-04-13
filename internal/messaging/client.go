// Package messaging 提供消息发送和接收功能
package messaging

import (
	"github.com/jason/agent-speaker/pkg/types"
)

// Client 消息客户端
type Client struct {
	identity *types.Identity
	encrypter types.Encrypter
}

// NewClient 创建消息客户端
func NewClient(identity *types.Identity, encrypter types.Encrypter) *Client {
	return &Client{
		identity:  identity,
		encrypter: encrypter,
	}
}

// Send 发送消息
func (c *Client) Send(to *types.Contact, content string, encrypt bool) error {
	// 简化实现
	return nil
}
