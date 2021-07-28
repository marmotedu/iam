// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package load

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	redis "github.com/go-redis/redis/v7"
	"github.com/marmotedu/component-base/pkg/json"

	"github.com/marmotedu/iam/pkg/log"
	"github.com/marmotedu/iam/pkg/storage"
)

// NotificationCommand defines a new notification type.
type NotificationCommand string

// Define Redis pub/sub events.
const (
	RedisPubSubChannel                      = "iam.cluster.notifications"
	NoticePolicyChanged NotificationCommand = "PolicyChanged"
	NoticeSecretChanged NotificationCommand = "SecretChanged"
)

// Notification is a type that encodes a message published to a pub sub channel (shared between implementations).
type Notification struct {
	Command       NotificationCommand `json:"command"`
	Payload       string              `json:"payload"`
	Signature     string              `json:"signature"`
	SignatureAlgo crypto.Hash         `json:"algorithm"`
}

// Sign sign Notification with SHA256 algorithm.
func (n *Notification) Sign() {
	n.SignatureAlgo = crypto.SHA256
	hash := sha256.Sum256([]byte(string(n.Command) + n.Payload))
	n.Signature = hex.EncodeToString(hash[:])
}

func handleRedisEvent(v interface{}, handled func(NotificationCommand), reloaded func()) {
	message, ok := v.(*redis.Message)
	if !ok {
		return
	}

	notif := Notification{}
	if err := json.Unmarshal([]byte(message.Payload), &notif); err != nil {
		log.Errorf("Unmarshalling message body failed, malformed: ", err)

		return
	}
	log.Infow("receive redis message", "command", notif.Command, "payload", message.Payload)

	switch notif.Command {
	case NoticePolicyChanged, NoticeSecretChanged:
		log.Info("Reloading secrets and policies")
		reloadQueue <- reloaded
	default:
		log.Warnf("Unknown notification command: %q", notif.Command)

		return
	}

	if handled != nil {
		// went through. all others shoul have returned early.
		handled(notif.Command)
	}
}

// RedisNotifier will use redis pub/sub channels to send notifications.
type RedisNotifier struct {
	store   *storage.RedisCluster
	channel string
}

// Notify will send a notification to a channel.
func (r *RedisNotifier) Notify(notif interface{}) bool {
	if n, ok := notif.(Notification); ok {
		n.Sign()
		notif = n
	}

	toSend, err := json.Marshal(notif)
	if err != nil {
		log.Errorf("Problem marshaling notification: %s", err.Error())

		return false
	}

	log.Debugf("Sending notification: %v", notif)

	if err := r.store.Publish(r.channel, string(toSend)); err != nil {
		if !errors.Is(err, storage.ErrRedisIsDown) {
			log.Errorf("Could not send notification: %s", err.Error())
		}

		return false
	}

	return true
}
