package auth

import (
	"context"
)

// Token token认证
type Token struct {
	AppID     string
	AppSecret string
}

// GetRequestMetadata 获取当前请求认证所需的元数据（metadata）
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": t.AppID, "app_secret": t.AppSecret}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return true
}
