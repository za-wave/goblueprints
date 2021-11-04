package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Title   string         `json:"title`
	Options []string       `json:"options"`
	Results map[string]int `json: "results,omitempty"`
}

func handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePollsGet(w, r)
		return
	case "POST":
		handlePollsPost(w, r)
		return
	case "DELETE":
		handlePollsDelete(w, r)
		return
	}
	// 未対応のhttp Method
	respondHTTPErr(w, r, http.StatusNotFound)
}

func handlePollsGet(w http.ResponseWriter, r *http.Request) {
	// respondHTTPErr(w, r, http.StatusInternalServerError, errors.New("未実装"))
	db := GetVar(r, "db").(*mgo.Database)
	c := db.C("polls")
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		//特定の調査項目の詳細
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// すべての調査項目リスト
		q = c.FindId(nil)
	}
	var result []*poll
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, &result)
}

func handlePollsPost(w http.ResponseWriter, r *http.Request) {
	// respondHTTPErr(w, r, http.StatusInternalServerError, errors.New("未実装"))
}

func handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	// respondHTTPErr(w, r, http.StatusInternalServerError, errors.New("未実装"))
}
