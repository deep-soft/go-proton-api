package server

import (
	"sync"

	"github.com/ProtonMail/go-proton-api"
)

func NewAuthCache() AuthCacher {
	return &authCache{
		info: make(map[string]proton.AuthInfo),
		auth: make(map[string]proton.Auth),
	}
}

type authCache struct {
	info map[string]proton.AuthInfo
	auth map[string]proton.Auth
	lock sync.RWMutex
}

func (c *authCache) GetAuthInfo(username string) (proton.AuthInfo, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	info, ok := c.info[username]

	return info, ok
}

func (c *authCache) SetAuthInfo(username string, info proton.AuthInfo) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.info[username] = info
}

func (c *authCache) GetAuth(username string) (proton.Auth, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	auth, ok := c.auth[username]

	return auth, ok
}

func (c *authCache) SetAuth(username string, auth proton.Auth) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.auth[username] = auth
}

func (c *authCache) FindAuthByUID(uid string) (string, proton.Auth) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for username, auth := range c.auth {
		if auth.UID == uid {
			return username, auth
		}
	}

	return "", proton.Auth{}
}

func (c *authCache) PopAuth() (string, proton.Auth) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for username, auth := range c.auth {
		delete(c.auth, username)
		delete(c.info, username)

		return username, auth
	}

	return "", proton.Auth{}
}
