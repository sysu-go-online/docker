package token

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"
)

// User 对用用户
type User struct {
	// token's ID
	ID string
	// 用于删除token
	Timer *time.Timer
}

// 储存根据时间遗弃的token
var PrevToken map[string]string = make(map[string]string)

// 储存token的map
var Token map[string]*User = make(map[string]*User)

// token更新时间为24h
var maxLiveTime time.Duration = time.Hour * 24

// 为token加的读写锁
var lock sync.RWMutex

// New 生成新的token，由于token24h后会修改，因此token由指针传入
func New(userID string, token *string, callback func()) {
	lock.Lock()
	prevtoken := *token
	_, has := Token[*token]
	delete(Token, *token)
	hash := md5.New()
	io.WriteString(hash, userID)
	io.WriteString(hash, time.Now().String())
	*token = fmt.Sprintf("%x", hash.Sum(nil))
	timer := time.AfterFunc(maxLiveTime, func() { New(userID, token, callback) })
	Token[*token] = &User{
		ID:    userID,
		Timer: timer,
	}
	if has {
		for k, v := range PrevToken {
			if v == prevtoken {
				delete(PrevToken, k)
			}
		}
		PrevToken[prevtoken] = *token
	}
	lock.Unlock()
	if callback != nil {
		callback()
	}
}

// Del 根据给定的token删除
func Del(token string) {
	lock.Lock()
	defer lock.Unlock()
	for k, v := range PrevToken {
		if v == token {
			delete(PrevToken, k)
		}
	}
	Token[token].Timer.Stop()
	delete(Token, token)
}

// GetUser 通过给定的token得到对应的user
func GetUser(token string) (*User, bool) {
	lock.RLock()
	defer lock.RUnlock()
	user, exisToken := Token[token]
	return user, exisToken
}
