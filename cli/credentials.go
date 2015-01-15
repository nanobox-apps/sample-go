package main

import (
  "os"
  "flag"
)

var hostname string
var key string
var id string
var bucketid string
var objectid string
var filename string
var out string

func init() {
	flag.StringVar(&hostname, "location", "", "userhost")
	flag.StringVar(&hostname, "host", "", "userhost")
	flag.StringVar(&key, "key", "", "userkey")
	flag.StringVar(&key, "user-key", "", "userkey")
	flag.StringVar(&id, "id", "", "userkey")
	flag.StringVar(&id, "user-id", "", "userkey")
	flag.StringVar(&bucketid, "bucketid", "", "userkey")
	flag.StringVar(&bucketid, "bucketname", "", "userkey")
	flag.StringVar(&objectid, "objectid", "", "userkey")
	flag.StringVar(&objectid, "objectalias", "", "userkey")
	flag.StringVar(&filename, "filename", "", "filename")
	flag.StringVar(&filename, "f", "", "filename")
	flag.StringVar(&out, "out", "", "output file")
	flag.StringVar(&out, "o", "", "output file")

	flag.Parse()
	host()
	userKey()
	userId()
	bucketId()
	objectId()
}

func host() {
	val := hostname
	if val == "" {
		val = os.Getenv("HOST")
	}
	if val == "" {
		val = os.Getenv("LOCATION")
	}
	hostname = val
}

func userKey() {
	val := key
	if val == "" {
		val = os.Getenv("USERKEY")
	}
	if val == "" {
		val = os.Getenv("KEY")
	}
	key = val
}

func userId() {
	val := id
	if val == "" {
		val = os.Getenv("USERID")
	}
	if val == "" {
		val = os.Getenv("ID")
	}
	id = val
}

func bucketId() {
	val := bucketid
	if val == "" {
		val = os.Getenv("BUCKETID")
	}
	if val == "" {
		val = os.Getenv("BUCKETNAME")
	}
	bucketid = val
}

func objectId() {
	val := objectid
	if val == "" {
		val = os.Getenv("OBJECTID")
	}
	if val == "" {
		val = os.Getenv("OBJECTALIAS")
	}
	objectid = val
}

