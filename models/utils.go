package models

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	salt1 = "sdfsfwef"
	salt2 = "XXSDsfdsdfsfwef"
)

func SecureID(prefix string, objectID primitive.ObjectID) string {
	now := time.Now()
	//hash:="sdf"
	data0 := []byte(fmt.Sprintf("%s???%d%s%s", salt2, now.Unix(), salt1, objectID.Hex()))
	data := []byte(fmt.Sprintf("%X%s???%s#%s-SQQAAAL%x%X", sha1.Sum(data0), salt1,
		prefix, objectID.Hex(), now.Unix(), sha512.Sum512(data0)))
	return fmt.Sprintf("%s#%s-%x-%x", prefix, objectID.Hex(), now.Unix(), md5.Sum(data))
}
