package lambda_service

import (
	"fmt"
	"testing"
)

type TestObj struct {
	SN         string
	Isp        string
	Provience  string
	LambdaName []string
}

func TestRadixIndex(t *testing.T) {
	t.Run("idx", func(t *testing.T) {
		t.Log("test begin")

		ispName := []string{"ct", "cnc", "cmnet"}
		provienceName := []string{"shangdong", "hebei", "guangdong", "hunan"}
		lambdaName := []string{"l1", "l2"}

		idx := NewRadixIndex(10)

		for i := 0; i < 10; i++ {
			o := TestObj{
				SN:         fmt.Sprintf("t%d", i),
				Isp:        ispName[i%len(ispName)],
				Provience:  provienceName[i%len(provienceName)],
				LambdaName: lambdaName,
			}

			idx.Insert(&o)
		}

		idx.CreateIndex(func(v interface{}) ([]string, error) {
			pathSlice := []string{}
			o := v.(*TestObj)
			for _, n := range o.LambdaName {
				path := n + "/" + o.Isp + "/" + o.Provience
				pathSlice = append(pathSlice, path)
			}

			return pathSlice, nil
		})

		objSliceI := idx.FindPrefix("l1")
		objSlice := []*TestObj{}
		for _, oI := range objSliceI {
			objSlice = append(objSlice, oI.(*TestObj))
		}

		t.Logf("objSlice: %d", len(objSlice))
		// for _, o := range objSlice {
		// 	t.Log(*o)
		// }

		if len(objSlice) != 10 {
			t.Errorf("objSlice len %d !=10", len(objSlice))
			return
		}

		for _, o := range objSlice {
			if o.SN == "t1" {
				if o.Isp != "cnc" || o.Provience != "hebei" {
					t.Errorf("%v not match t1 hebei", o)
				}
			}
			if o.SN == "t3" {
				if o.Isp != "ct" || o.Provience != "hunan" {
					t.Errorf("%v not match t1 hebei", o)
				}
			}
			if o.SN == "t8" {
				if o.Isp != "cmnet" || o.Provience != "shangdong" {
					t.Errorf("%v not match t1 hebei", o)
				}
			}
		}

	})
}
