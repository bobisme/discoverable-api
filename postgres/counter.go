package postgres

//go:generate msgp
import (
	"bytes"
	"database/sql"
	"encoding/binary"

	"github.com/tinylib/msgp/msgp"
)

type CacheKey struct {
	Stmt string        `msg:"stmt"`
	Args []interface{} `msg:"args"`
}

type cache interface {
	Get([]byte) ([]byte, error)
	Set(key, val []byte, ttl int) error
}

type Counter struct {
	cache cache
	db    *sql.DB
}

func (c *Counter) key(stmt string, args ...interface{}) []byte {
	buf := new(bytes.Buffer)
	err := msgp.Encode(buf, &CacheKey{stmt, []interface{}{}})
	fatalIfErr(err)
	return buf.Bytes()
}

func (c *Counter) set(key []byte, count int) {
	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(count))

	fatalIfErr(c.cache.Set(key, countBytes, 60))
}

func (c *Counter) get(key []byte) (int, error) {
	out, err := c.cache.Get(key)
	if err != nil {
		return 0, err
	}
	count := binary.LittleEndian.Uint64(out)
	return int(count), nil
}

func (c *Counter) Count(stmt string, args ...interface{}) (int, error) {
	key := c.key(stmt, args...)
	if count, err := c.get(key); err == nil {
		return count, nil
	}
	var count int
	err := c.db.QueryRow(stmt, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	c.set(key, count)
	return count, nil
}
