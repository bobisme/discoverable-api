package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bobisme/discoverable-api/postgres"
	"github.com/bobisme/discoverable-api/storage"
	"github.com/coocood/freecache"
	_ "github.com/lib/pq"
	metrics "github.com/rcrowley/go-metrics"
)

const fakeRowCount = 10
const limit = 1000

func fatalIfErr(possibleErrors ...interface{}) {
	if err, ok := possibleErrors[len(possibleErrors)-1].(error); ok {
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createTable(db *sql.DB) {
	mustExec := func(stmt string, args ...interface{}) {
		fatalIfErr(db.Exec(stmt, args...))
	}
	_, err := fmt.Fprintf(
		os.Stderr, "creating things table, %d rows...", fakeRowCount)
	fatalIfErr(err)
	mustExec(`
		CREATE TEMPORARY TABLE things AS
			SELECT
				generate_series(1, $1) AS n,
				substr(lower(md5(random()::text)), 1, 16) AS id,
				substr(
					concat(md5(random()::text), md5(random()::text)),
					1, (random() * 64)::integer + 1
				) AS description`, fakeRowCount)
	mustExec(`CREATE INDEX things_id_idx ON things USING btree (id)`)
	mustExec(`
		CREATE TEMPORARY TABLE widgets (
			id VARCHAR(22) PRIMARY KEY,
			name VARCHAR(128),
			created_at DATETIME NOT NULL,
			deleted_at DATETIME
		)`)
	mustExec(`
		CREATE INDEX created_at__id ON widgets
		USING btree (created_at, id)`)
	mustExec(`
		CREATE TEMPORARY TABLE widget_things (
			widget_id VARCHAR(22) REFERENCES widgets (id)
				ON UPDATE CASCADE ON DELETE CASCADE,
			thing_id VARCHAR(16) REFERENCES things (id)
				ON UPDATE CASCADE
			CONSTRAINT widget_thing_pkey PRIMARY KEY (widget_id, thing_id)
		)`)
	_, err = fmt.Fprintln(os.Stderr, "...done")
	fatalIfErr(err)
}

func main() {
	connStr := "user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	fatalIfErr(err)

	createTable(db)

	cache := freecache.NewCache(10 * 1 << 20)

	log.Println("selecting in groups of", limit)

	var thingStore storage.ThingStorage = postgres.NewThingStorage(
		db, cache)
	// var thing *models.Thing = thingStore.Get("123")
	// log.Println(thing)

	t := metrics.NewTimer()
	fatalIfErr(metrics.Register("page", t))
	ttotal := metrics.NewTimer()

	go metrics.LogScaled(
		metrics.DefaultRegistry, 1*time.Second, time.Millisecond,
		log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	timeList := func(cursor string) storage.ThingPage {
		// log.Println("cursor:", cursor)
		var things storage.ThingPage
		t.Time(func() {
			things = thingStore.List(cursor, limit)
		})
		out, _ := json.Marshal(things)
		fmt.Println(string(out))
		return things
	}

	var things storage.ThingPage
	ttotal.Time(func() {
		timeList("")
		things = timeList("")
		nextCursor := things.NextCursor
		for nextCursor != "" {
			things = timeList(nextCursor)
			nextCursor = things.NextCursor
		}
	})
	// log.Println("GOING BACKWARDS")
	// nextCursor := things.PrevCursor
	// for nextCursor != "" {
	// 	things = timeList(nextCursor)
	// 	nextCursor = things.PrevCursor
	// }

	// fatalIfErr(metrics.Register("total", ttotal))
	time.Sleep(1 * time.Second)
}
