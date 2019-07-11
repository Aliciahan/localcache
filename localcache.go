package localcache

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/gregjones/httpcache/diskcache"
	"log"
	"time"
)

type LocalCache struct {
	dirPath string
}

func (l LocalCache) GetDirPath () string {
	return l.dirPath
}

func (l LocalCache) SetDirPath (dirPath string) {
	l.dirPath = dirPath
}

func NewLocalCache (dirPath string) LocalCache {
	return  LocalCache{
		dirPath: dirPath,
	}
}


type dataWithTimeStamp struct {
	TimeAdded 	int64
	Data 		[]byte
}

func int64Tobytes(i int64) []byte {
	b := make([]byte,8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func bytesToInt64(b []byte ) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

func getBytesFromStruct(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getInterface(bts []byte, data interface{})error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func (l LocalCache) SaveData (key string, d []byte){
	c := diskcache.New(l.dirPath)
	dd, err:= getBytesFromStruct(dataWithTimeStamp{
		TimeAdded: time.Now().UnixNano(),
		Data: d,
	})
	if err != nil {
		log.Fatal("Encoding Error: ", err)
	}
	c.Set(key,dd)
}

func (l LocalCache) GetData (key string, ttlInSecond int) ([]byte, bool){
	var ttl int64
	ttl = 1000000000*int64(ttlInSecond)
	c := diskcache.New(l.dirPath)
	dataBody,existe := c.Get(key)
	if existe == false {
		//log.Println("Key does not existe ... ")
		return nil, false
	}
	var data dataWithTimeStamp

	err := getInterface(dataBody,&data)
	if err != nil {
		log.Fatal("interface",err)
	}

	if time.Now().UnixNano()-data.TimeAdded <= ttl {
		return data.Data, true
	}else {
		//log.Println("Key expired ... ")
		return nil, false
	}
}

func (l LocalCache) GetDataWithStatus (key string, ttlInSecond int) ([]byte, int){
	/*
	return
		0: The data exists and within TTL
		1: The data doesn't exist
		2: The data exists but not in a valid TTL
	 */
	var ttl int64
	ttl = 1000000000*int64(ttlInSecond)
	c := diskcache.New(l.dirPath)
	dataBody,existe := c.Get(key)
	if existe == false {
		//log.Println("Key does not existe ... ")
		return nil, 1
	}
	var data dataWithTimeStamp

	err := getInterface(dataBody,&data)
	if err != nil {
		log.Fatal("interface",err)
	}

	if time.Now().UnixNano()-data.TimeAdded <= ttl {
		return data.Data, 0
	}else {
		//log.Println("Key expired ... ")
		return nil, 2
	}
}
