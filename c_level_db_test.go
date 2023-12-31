// +build gcc

package db

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/nio-net/common"
)

func BenchmarkRandomReadsWrites2(b *testing.B) {
	b.StopTimer()

	numItems := int64(1000000)
	internal := map[int64]int64{}
	for i := 0; i < int(numItems); i++ {
		internal[int64(i)] = int64(0)
	}
	db, err := NewCLevelDB(Fmt("test_%x", RandStr(12)), "")
	if err != nil {
		b.Fatal(err.Error())
		return
	}

	fmt.Println("ok, starting")
	b.StartTimer()

	for i := 0; i < b.N; i++ {

		{
			idx := (int64(RandInt()) % numItems)
			internal[idx] += 1
			val := internal[idx]
			idxBytes := int642Bytes(int64(idx))
			valBytes := int642Bytes(int64(val))

			db.Set(
				idxBytes,
				valBytes,
			)
		}

		{
			idx := (int64(RandInt()) % numItems)
			val := internal[idx]
			idxBytes := int642Bytes(int64(idx))
			valBytes := db.Get(idxBytes)

			if val == 0 {
				if !bytes.Equal(valBytes, nil) {
					b.Errorf("Expected %X for %v, got %X",
						nil, idx, valBytes)
					break
				}
			} else {
				if len(valBytes) != 8 {
					b.Errorf("Expected length 8 for %v, got %X",
						idx, valBytes)
					break
				}
				valGot := bytes2Int64(valBytes)
				if val != valGot {
					b.Errorf("Expected %v for %v, got %v",
						val, idx, valGot)
					break
				}
			}
		}
	}

	db.Close()
}
