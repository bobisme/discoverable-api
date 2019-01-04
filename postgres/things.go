package postgres

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"log"

	"github.com/bobisme/discoverable-api/models"
	"github.com/bobisme/discoverable-api/storage"
)

const defaultLimit = 4

func fatalIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type countCache struct {
	cache cache
}

func (c countCache) set(stmt string, count int) {
	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(count))
	fatalIfErr(c.cache.Set([]byte(stmt), countBytes, 60))
}
func (c countCache) get(stmt string) (int, error) {
	out, err := c.cache.Get([]byte(stmt))
	if err != nil {
		return 0, err
	}
	count := binary.LittleEndian.Uint64(out)
	return int(count), nil
}

func encodeCursor(thingId string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(thingId))
}

func decodeCursor(cursor string) string {
	id, err := base64.RawURLEncoding.DecodeString(cursor)
	fatalIfErr(err)
	return string(id)
}

type ThingStorage struct {
	db         *sql.DB
	countCache countCache
}

func NewThingStorage(db *sql.DB, cache cache) *ThingStorage {
	return &ThingStorage{
		db:         db,
		countCache: countCache{cache},
	}
}

func (t *ThingStorage) Get(id string) *models.Thing {
	return &models.Thing{Id: "123", Description: "hi there"}
}

func (t *ThingStorage) List(cursor string, limit int) storage.ThingPage {
	// SELECT * FROM medley WHERE n > 99999 ORDER BY n ASC LIMIT 10;
	// SELECT * FROM medley WHERE n < 99999 ORDER BY n DESC OFFSET 9 LIMIT 1;
	cursorId := decodeCursor(cursor)
	lister := &thingLister{
		db:         t.db,
		countCache: t.countCache,
		cursorId:   cursorId,
		limit:      limit,
	}
	if lister.limit <= 0 {
		lister.limit = defaultLimit
	}
	return lister.list()
}

type thingLister struct {
	db         *sql.DB
	countCache countCache
	cursorId   string
	limit      int
}

func (l *thingLister) list() storage.ThingPage {
	things := l.queryPage()
	next := l.nextCursor(things)
	if len(things) > l.limit {
		things = things[:l.limit]
	}
	return storage.ThingPage{
		Total:      l.countAll(),
		PrevCursor: l.prevCursor(),
		NextCursor: next,
		Items:      things,
	}
}

func (l *thingLister) countAll() int {
	var count int

	stmt := `SELECT COUNT(*) FROM things`
	{
		if count, err := l.countCache.get(stmt); err == nil {
			return count
		}
	}
	// {
	// 	var countEstimate float64
	// 	err := l.db.QueryRow(`
	// 		SELECT
	// 		  (reltuples/relpages) * (
	// 			pg_relation_size('things') /
	// 			(current_setting('block_size')::integer)
	// 		  )
	// 		  FROM pg_class where relname = 'things'
	// 	`).Scan(&countEstimate)
	// 	count = int(math.Round(countEstimate))
	// 	fatalIfErr(err)
	// }
	// if count < 1000 {
	fatalIfErr(l.db.QueryRow(stmt).Scan(&count))
	// }
	l.countCache.set(stmt, count)
	return count
}

func (l *thingLister) prevCursor() string {
	var prevCursorId string
	err := l.db.QueryRow(
		`SELECT id FROM things
			WHERE id < $1
			ORDER BY id DESC
			OFFSET $2
			LIMIT 1`, l.cursorId, l.limit-1).
		Scan(&prevCursorId)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}
	return encodeCursor(prevCursorId)
}

func (l *thingLister) queryRows() (*sql.Rows, error) {
	return l.db.Query(
		`SELECT id, description
				FROM things
				WHERE id >= $1
				ORDER BY id ASC
				LIMIT $2`,
		l.cursorId, l.limit+1)
}

func (l *thingLister) queryPage() []*models.Thing {
	rows, err := l.queryRows()
	defer rows.Close()
	fatalIfErr(err)
	var things []*models.Thing
	for rows.Next() {
		thing := new(models.Thing)
		err := rows.Scan(&thing.Id, &thing.Description)
		fatalIfErr(err)
		things = append(things, thing)
	}
	fatalIfErr(rows.Err())
	return things
}

func (l *thingLister) nextCursor(things []*models.Thing) string {
	if len(things) > l.limit {
		return encodeCursor(things[l.limit].Id)
	}
	return ""
}
