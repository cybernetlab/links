package models

import (
	"context"
	"crypto/md5"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Link struct {
	UserID    uint      `json:"user_id"`
	Link      string    `json:"link"`
	Count     uint      `json:"count"`
	CreatedAt time.Time `json:"time"`
	Hash      string    `json:"-"`
}

func CreateLink(uid uint, link string, exp time.Duration) (*Link, error) {
	hash := createHash(fmt.Sprintf("%v:%v", uid, link))
	key := fmt.Sprintf("%v:%v", hash, uid)
	l := &Link{UserID: uid, Link: link, Count: 0, CreatedAt: time.Now(), Hash: hash}
	value, _ := json.Marshal(l)
	err := KV.SetEx(context.Background(), key, value, exp).Err()
	return l, err
}

func DeleteLink(uid uint, hash string) {
	key := fmt.Sprintf("%v:%v", hash, uid)
	KV.Del(context.Background(), key)
}

func UpdateLink(link *Link) {
	if link == nil {
		return
	}
	key := fmt.Sprintf("%v:%v", link.Hash, link.UserID)
	value, _ := json.Marshal(link)
	KV.Set(context.Background(), key, value, redis.KeepTTL)
}

func LinkByHash(hash string) *Link {
	ctx := context.Background()
	pattern := fmt.Sprintf("%v:*", hash)

	var keys []string
	var value string
	var err error
	keys, _, err = KV.Scan(ctx, 0, pattern, 1).Result()
	if err != nil || len(keys) != 1 {
		return nil
	}
	value, err = KV.Get(ctx, keys[0]).Result()
	if err != nil {
		return nil
	}
	return deserializeLink(keys[0], value)
}

func LinksByUserID(uid uint) []Link {
	result := make([]Link, 0)
	ctx := context.Background()
	pattern := fmt.Sprintf("*:%v", uid)

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = KV.Scan(ctx, cursor, pattern, 50).Result()
		if err != nil {
			break
		}

		links, err := KV.MGet(ctx, keys...).Result()
		if err != nil {
			break
		}
		for i, encoded := range links {
			link := deserializeLink(keys[i], encoded)
			if link == nil {
				continue
			}
			result = append(result, *link)
		}

		if cursor == 0 {
			break
		}
	}

	return result
}

func createHash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	bytes := h.Sum(nil)
	encoded := base32.StdEncoding.EncodeToString(bytes)
	return strings.TrimRight(encoded, "=")
}

func deserializeLink(key string, encoded interface{}) *Link {
	var link Link
	str, ok := encoded.(string)
	if !ok {
		return nil
	}
	if err := json.Unmarshal([]byte(str), &link); err != nil {
		return nil
	}
	parts := strings.SplitN(key, ":", 2)
	link.Hash = parts[0]
	return &link
}
