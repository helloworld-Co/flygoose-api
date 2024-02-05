package caching

import (
	"fmt"
	"github.com/fanjindong/go-cache"
	"sync"
	"time"
)

var (
	baseAdminCacheKey = "admin-"
	keyAdSmsCode      = baseAdminCacheKey + "sms-code-%s"    //短信验证码
	keyAdAdminToken   = baseAdminCacheKey + "admin-token-%s" //访问token
)

var (
	adminCacheEngine *AdminCacheEngine
	lockH            sync.Mutex
)

type AdminCacheEngine struct {
	c cache.ICache
}

func newAdminCache() *AdminCacheEngine {
	// go-cache会定时清理过期的缓存对象，默认间隔是5秒
	return &AdminCacheEngine{c: cache.NewMemCache(cache.WithClearInterval(5 * time.Second))}
}

// GetAdminCache 单例，管理admin服务缓存
func GetAdminCache() *AdminCacheEngine {
	if adminCacheEngine != nil {
		return adminCacheEngine
	}

	lockH.Lock()
	defer lockH.Unlock()

	if adminCacheEngine != nil {
		return adminCacheEngine
	}

	adminCacheEngine = newAdminCache()
	return adminCacheEngine
}

func (rev *AdminCacheEngine) GetUid(token string) int64 {
	if v, ok := rev.c.Get(token); ok {
		if uid, ok := v.(int64); ok {
			return uid
		}
	}
	return 0
}

func (rev *AdminCacheEngine) SetUid(token string, uid int64) bool {
	return rev.c.Set(token, uid)
}

func (rev *AdminCacheEngine) DeleteUid(token string) {
	rev.c.Del(token)
}

func (rev *AdminCacheEngine) SetSmsCode(phone string, smsCode string) bool {
	return rev.c.Set(fmt.Sprintf(keyAdSmsCode, phone), smsCode, cache.WithEx(60*time.Second)) //默认有效期60s
}

func (rev *AdminCacheEngine) GetSmsCode(phone string) string {
	if s, ok := rev.c.Get(fmt.Sprintf(keyAdSmsCode, phone)); ok {
		return s.(string)
	}
	return ""
}

func (rev *AdminCacheEngine) SetAdminToken(phone string, token string) bool {
	return rev.c.Set(fmt.Sprintf(keyAdAdminToken, phone), token, cache.WithEx(24*time.Hour)) //默认有效期24小时
}

func (rev *AdminCacheEngine) GetAdminToken(phone string) string {
	if s, ok := rev.c.Get(fmt.Sprintf(keyAdAdminToken, phone)); ok {
		return s.(string)
	}
	return ""
}
