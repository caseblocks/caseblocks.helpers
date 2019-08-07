package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func jid() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s\n", hex.EncodeToString(b)), nil
	}
}

type SidekiqJob struct {
	Retry      bool                     `json:"retry"`
	Queue      string                   `json:"queue"`
	Class      string                   `json:"class"`
	Args       []map[string]interface{} `json:"args"`
	Jid        string                   `json:"jid"`
	EnqueuedAt int64                    `json:"enqueued_at"`
}

func LauchSidekiqJob(queue string, class string, params []map[string]interface{}) error {

	c, err := redis.Dial("tcp", FindConnString())
	if err != nil {
		return err
	}
	defer c.Close()

	jid, _ := jid()
	skj := SidekiqJob{true, queue, class, params, jid, time.Now().UTC().Unix()}

	if skjJsonBytes, err := json.Marshal(skj); err != nil {
		return err
	} else {
		_, err := c.Do("RPUSH", "queue:system", fmt.Sprintf("%s", skjJsonBytes))
		return err
	}
	return nil
}
