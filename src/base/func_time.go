package base

import (
	"time"
	//"fmt"
	//"strconv"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// func time_now (expr *CXExpression, call *CXCall) error {
// 	t := time.Now().Format("20060102150405")
// 	timestamp, err := strconv.ParseInt(t, 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	sT := encoder.Serialize(timestamp)
// 	assignOutput(0, sT, "i64", expr, call)
// 	return nil
// }

func time_now (expr *CXExpression, call *CXCall) error {
	// t := time.Now()
	// ts := t.Format("20060102150405")
	// sTime := fmt.Sprintf("%s.%d", ts, t.Nanosecond())
	// time, err := strconv.ParseFloat(sTime, 64)

	// if err != nil {
	// 	return err
	// }

	// sT := encoder.Serialize(time)
	// assignOutput(0, sT, "f64", expr, call)
	// return nil

	t := float64(time.Now().UnixNano()) / 1000000000 // to seconds
	sT := encoder.Serialize(t)
	assignOutput(0, sT, "f64", expr, call)
	return nil
}
