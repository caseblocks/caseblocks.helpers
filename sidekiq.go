package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
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

	// logger := cb.Logger()
	// db := cb.NewSqlConnection(os.Getenv("MYSQL_CONN"), logger)

	// bucket := cb.Bucket{}
	// rows, err := db.Queryx("select id, name, case_type_id, kpi, last_checked_membership_at, last_checked_tripping_at from case_blocks_buckets;")
	// for rows.Next() {
	//  err := rows.StructScan(&bucket)
	//  cb.PanicIf(err, logger)
	//  logger.Printf("%#v\n", bucket)
	// }

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer c.Close()

	jid, _ := jid()
	skj := SidekiqJob{true, queue, class, params, jid, time.Now().UTC().Unix()}

	if skjJsonBytes, err := json.Marshal(skj); err != nil {
		return err
	} else {
		c.Do("RPUSH", "queue:system", fmt.Sprintf("%s", skjJsonBytes))
	}
	return nil
}
