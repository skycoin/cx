package base

// import (
// 	"time"
// 	"github.com/skycoin/skycoin/src/cipher/encoder"
// )

// func time_Unix (expr *CXExpression, call *CXCall) error {
// 	t := time.Now().Unix()
// 	sT := encoder.Serialize(t)
// 	assignOutput(0, sT, "i64", expr, call)
// 	return nil
// }

// func time_UnixMilli (expr *CXExpression, call *CXCall) error {
// 	//t := float64(time.Now().UnixNano()) / 1000000000
// 	t := time.Now().UnixNano() * 1000000
// 	sT := encoder.Serialize(t)
// 	assignOutput(0, sT, "i64", expr, call)
// 	return nil
// }

// func time_UnixNano (expr *CXExpression, call *CXCall) error {
// 	t := time.Now().Unix()
// 	sT := encoder.Serialize(t)
// 	assignOutput(0, sT, "i64", expr, call)
// 	return nil
// }
