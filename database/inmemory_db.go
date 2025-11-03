package database

import "sync"

var db *sync.Map

func InitDataBase() {
	db = &sync.Map{}
}

func Store(key string, value interface{}) {
	db.Store(key, value)
}

func Get(key string) (interface{}, bool) {
	return db.Load(key)
}

func Delete(key string) {
	db.Delete(key)
}
