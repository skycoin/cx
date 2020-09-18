//+build cxstrat

package model

import (
	"errors"
	"fmt"
	"time"

	. "github.com/SkycoinProject/cx/cx"
	"github.com/boltdb/bolt"
)

type CXDB struct {
	b *bolt.DB
}

var cxdb *CXDB

func CXDBOpen() {
	cxdb = new(CXDB)
	b_, err := bolt.Open("cxdb.db", 0644, &bolt.Options{
		Timeout: time.Millisecond * 5000,
	})

	cxdb.b = b_

	if err != nil {
		panic("CXDB: cannot open!")
	}
}

func newBucket(prgm *CXProgram) {
	if cxdb == nil || cxdb.b == nil {
		CXDBOpen()
	}

	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	buckname := ReadStr(fp, expr.Inputs[0])

	err := cxdb.b.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(buckname))
		return
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), err == nil)
}

func fetch(prgm *CXProgram) {
	if cxdb == nil || cxdb.b == nil {
		CXDBOpen()
	}

	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	buckname := []byte(ReadStr(fp, expr.Inputs[0]))
	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[1])
	query := GetSliceData(inputSliceOffset, 1)

	_ = cxdb.b.View(func(tx *bolt.Tx) (err error) {
		wbkt := tx.Bucket(buckname)
		if wbkt != nil {
			val := wbkt.Get(query)
			if val != nil {
				outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
				outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
				outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(val)), 1))
				copy(GetSliceData(outputSliceOffset, 1), val)
				WriteI32(outputSlicePointer, outputSliceOffset)
			} else {
				WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(0))
			}
		} else {
			WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(0))
		}
		return
	})
}

func store(prgm *CXProgram) {
	if cxdb == nil || cxdb.b == nil {
		CXDBOpen()
	}

	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	buckname := []byte(ReadStr(fp, expr.Inputs[0]))
	inputSliceOffset1 := GetSliceOffset(fp, expr.Inputs[1])
	key := GetSliceData(inputSliceOffset1, 1)
	inputSliceOffset2 := GetSliceOffset(fp, expr.Inputs[2])
	val := GetSliceData(inputSliceOffset2, 1)

	err := cxdb.b.Update(func(tx *bolt.Tx) (err error) {
		wbkt := tx.Bucket(buckname)
		if wbkt == nil {
			return errors.New("dummy")
		}

		err = wbkt.Put(key, val)
		return
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), err == nil)
}
