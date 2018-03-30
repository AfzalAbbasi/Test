package main

import (
	//"github.com/fjl/go-couchdb"

	"github.com/couchbase/gocb"

	"fmt"
)

const (
	Host = ""
)

func main() {
	cluster, _ := gocb.Connect(Host)
	bucket, _ := cluster.OpenBucket("account/08/b8/0907389f7c529cfad507f12691ed", "")
	fmt.Println(bucket.Name())

}
