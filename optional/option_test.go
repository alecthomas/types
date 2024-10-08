package optional_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
	. "github.com/alecthomas/types/optional"
	_ "modernc.org/sqlite" // Register SQLite driver.
)

func TestOptionFrom(t *testing.T) {
	yes := func() (string, bool) { return "yes", true }
	no := func() (string, bool) { return "", false }

	o := From(yes())
	assert.Equal(t, Some("yes"), o)
	o = From(no())
	assert.Equal(t, None[string](), o)
}

func TestOptionGet(t *testing.T) {
	o := Some(1)
	v, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	o = None[int]()
	_, ok = o.Get()
	assert.False(t, ok)
}

func TestOptionMarshalJSON(t *testing.T) {
	o := Some(1)
	b, err := o.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "1", string(b))

	o = None[int]()
	b, err = o.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestOptionUnmarshalJSON(t *testing.T) {
	o := Option[int]{}
	err := json.Unmarshal([]byte("1"), &o)
	assert.NoError(t, err)
	b, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, b)
}

func TestOptionMarshalJSONOmitEmpty(t *testing.T) {
	s := struct {
		Value Option[int] `json:"value,omitempty"`
	}{}

	b, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":null}`, string(b))
}

func TestOptionString(t *testing.T) {
	o := Some(1)
	assert.Equal(t, "1", o.String())

	o = None[int]()
	assert.Equal(t, "None", o.String())
}

func TestOptionGoString(t *testing.T) {
	o := Some(1)
	assert.Equal(t, "Some[int](1)", o.GoString())

	o = None[int]()
	assert.Equal(t, "None[int]()", o.GoString())
}

func TestOptionSQL(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, value INTEGER);`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO test (id, value) VALUES (1, 1);`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO test (id, value) VALUES (2, NULL);`)
	assert.NoError(t, err)

	var option Option[int64]
	rows, err := db.Query("SELECT value FROM test WHERE id = 1;")
	assert.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&option)
		assert.NoError(t, err)
	}
	err = rows.Err()
	assert.NoError(t, err)
	assert.Equal(t, Some(int64(1)), option)

	rows, err = db.Query("SELECT value FROM test WHERE id = 2;")
	assert.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&option)
		assert.NoError(t, err)
	}
	err = rows.Err()
	assert.NoError(t, err)
	assert.Equal(t, None[int64](), option)
}

func TestOptionZero(t *testing.T) {
	assert.Equal(t, None[error](), Zero((error)(nil)))
	assert.Equal(t, None[string](), Zero(""))
}

func TestOptionNil(t *testing.T) {
	assert.Panics(t, func() {
		var str string
		assert.Equal(t, None[string](), Nil(str))
	})

	assert.Equal(t, None[error](), Nil((error)(nil)))
}
