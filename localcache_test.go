package localcache

import (
	"fmt"
	"log"
	"testing"
)

func TestStuff(t *testing.T) {

	//go test -test.run=TestStuff

	fmt.Println("hello")
	//lc := LocalCache{"/tmp/lc-test"}

	a := int64(123345)
	aByte := int64Tobytes(a)

	b := dataWithTimeStamp{
		TimeAdded: 1234,
		Data: aByte,
	}

	bb, err := getBytesFromStruct(b)
	if err != nil {
		log.Fatal("get",err)
	}

	var cc dataWithTimeStamp
	err = getInterface(bb, &cc)
	if err != nil {
		log.Fatal("interface",err)
	}
	log.Println("cc: ", cc.TimeAdded)
	log.Println("cc: ", cc.Data)

}