package jsontime

import (
	"bytes"
	"testing"
	"time"
)

var json = ConfigWithCustomTimeFormat

type Book struct {
	Id          int        `json:"id"`
	PublishedAt *time.Time `json:"published_at" time_format:"sql_date" time_utc:"true"`
	UpdatedAt   *time.Time `json:"updated_at" time_format:"sql_date" time_utc:"true"`
	CreatedAt   time.Time  `json:"created_at" time_format:"2006-01-02 15:04:05"`
}

func TestMarshalFormat(t *testing.T) {
	t2018 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	book := Book{
		Id:        1,
		UpdatedAt: &t2018,
		CreatedAt: t2018,
	}

	if b, err := json.Marshal(book); err != nil {
		t.Error(err)
	} else if string(b) != `{"id":1,"published_at":null,"updated_at":"2018-01-01","created_at":"2018-01-01 08:00:00"}` {
		t.Errorf("got:%s\n", b)
	}

}

func TestUnmarshalFormat(t *testing.T) {
	t2018utc := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	t2018 := time.Date(2018, 1, 1, 8, 0, 0, 0, time.Local)
	b := []byte(`{"id":1,"updated_at":"2018-01-01","created_at":"2018-01-01 08:00:00"}`)

	book := Book{}
	if err := json.Unmarshal(b, &book); err != nil {
		t.Error(err)
	}

	if book.Id != 1 || book.CreatedAt != t2018 ||
		book.UpdatedAt == nil || *book.UpdatedAt != t2018utc ||
		book.PublishedAt != nil {
		t.Errorf("got:%v", book)
	}
}

func TestDecoderFormat(t *testing.T) {
	t2018utc := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	t2018 := time.Date(2018, 1, 1, 8, 0, 0, 0, time.Local)
	b := []byte(`{"id":1,"updated_at":"2018-01-01","created_at":"2018-01-01 08:00:00"}`)

	decoder := json.NewDecoder(bytes.NewReader(b))
	book := Book{}
	if err := decoder.Decode(&book); err != nil {
		t.Error(err)
	}

	if book.Id != 1 || book.CreatedAt != t2018 ||
		book.UpdatedAt == nil || *book.UpdatedAt != t2018utc ||
		book.PublishedAt != nil {
		t.Errorf("got:%v", book)
	}

	if book.CreatedAt.Format("2006-01-02 15:04:05") != "2018-01-01 08:00:00" {
		t.Errorf("got:%v", book)
	}
}

type User struct {
	Id        int        `json:"id"`
	UpdatedAt *time.Time `json:"updated_at" time_format:"sql_datetime" time_location:"Local"`
	CreatedAt time.Time  `json:"created_at" time_format:"sql_datetime" time_location:"Local"`
}

func TestLocale(t *testing.T) {
	user := User{
		Id:        0,
		UpdatedAt: nil,
		CreatedAt: time.Date(0, 1, 1, 0, 0, 0, 0, time.Local),
	}

	b, err := json.Marshal(user)
	if err != nil {
		t.Error(err.Error())
	}

	if string(b) != `{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}` {
		t.Errorf("got: %s", b)
	}
}

func TestUnMarshalZero(t *testing.T) {
	user := User{}
	jsonBytes := []byte(`{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}`)

	err := json.Unmarshal(jsonBytes, &user)
	if err != nil {
		t.Error(err.Error())
	}

	b, err := json.Marshal(user)
	if err != nil {
		t.Error(err.Error())
	}

	if string(b) != `{"id":0,"updated_at":null,"created_at":"0000-00-00 00:00:00"}` {
		t.Errorf("got: %s", b)
	}
}

func TestTagDefault(t *testing.T) {
	str := `{"dt":"1212.234"}`

	p := &struct {
		Dt time.Time `time_format:"2020-10-00 01:02:01,default"`
	}{}

	{
		err := json.Unmarshal([]byte(str), p)
		if err != nil {
			t.Error("time parse return error")
		}
	}
}
