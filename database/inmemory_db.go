package database

import "sync"

var db *sync.Map
var keylock *sync.Map

func InitDataBase() {
	db = &sync.Map{}
	keylock = &sync.Map{}
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

func GetKeyLock(key string) *sync.Mutex {
	lock, _ := keylock.LoadOrStore(key, &sync.Mutex{})
	return lock.(*sync.Mutex)
}
