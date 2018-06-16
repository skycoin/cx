package main

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"io/ioutil"
	"os"
	"time"
)

type Point struct {
	x int32
	y int32
}

type Arrays struct {
	ints   []int32
	floats []float32
}

type Structs struct {
	points []Point
	arrays []Arrays
}

func main() {

	fileName := "serial"

	/*
	   Point
	*/

	testPoint := Point{
		x: 10,
		y: 20,
	}
	_ = testPoint

	/*
	   Arrays
	*/

	// ints := make([]int32, 200000, 200000)
	// floats := make([]float32, 200000, 200000)
	testArrays := Arrays{
		ints:   []int32{1, 2, 3},
		floats: []float32{1.0, 2.0, 3.0},
	}
	// testArrays := Arrays{
	// 	ints: ints,
	// 	floats: floats,
	// }

	/*
	   Structs
	*/

	point1 := Point{
		x: 10,
		y: 20,
	}

	point2 := Point{
		x: 20,
		y: 30,
	}

	arrays1 := Arrays{
		ints:   []int32{1, 2, 3},
		floats: []float32{1.0, 2.0, 3.0},
	}

	var points []Point
	points = append(points, point1)
	points = append(points, point2)

	var arrays []Arrays
	arrays = append(arrays, arrays1)

	var testStructs Structs
	testStructs.points = points
	testStructs.arrays = arrays

	// /*
	//     Point
	// */

	// // serialization only
	// start := time.Now()
	// for c := 0; c < 100000; c++ {
	// 	encoder.Serialize(testPoint)
	// }
	// duration := time.Since(start)
	// fmt.Println("s\t", duration)

	// // serialization + write
	// start = time.Now()
	// for c := 0; c < 100000; c++ {
	// 	//byts := encoder.Serialize(testPoint)
	// 	ioutil.WriteFile(fileName, encoder.Serialize(testPoint), os.FileMode(644))
	// }
	// duration = time.Since(start)
	// fmt.Println("s+w\t", duration)

	// // write only
	// byts := encoder.Serialize(testPoint)
	// start = time.Now()
	// for c := 0; c < 100000; c++ {
	// 	ioutil.WriteFile(fileName, byts, os.FileMode(644))
	// }
	// duration = time.Since(start)
	// fmt.Println("w\t", duration)

	/*
	   Arrays
	*/

	// serialization only
	start := time.Now()
	for c := 0; c < 100000; c++ {
		encoder.Serialize(testArrays)
	}
	duration := time.Since(start)
	fmt.Println("s\t", duration)

	// serialization + write
	start = time.Now()
	for c := 0; c < 100; c++ {
		//byts := encoder.Serialize(testArrays)
		ioutil.WriteFile(fileName, encoder.Serialize(testArrays), os.FileMode(644))
	}
	duration = time.Since(start)
	fmt.Println("s+w\t", duration)

	// write only
	byts := encoder.Serialize(testArrays)
	start = time.Now()
	for c := 0; c < 100; c++ {
		ioutil.WriteFile(fileName, byts, os.FileMode(644))
	}
	duration = time.Since(start)
	fmt.Println("w\t", duration)

	// /*
	//     Structs
	// */

	// // serialization only
	// start = time.Now()
	// for c := 0; c < 100000; c++ {
	// 	encoder.Serialize(testStructs)
	// }
	// duration = time.Since(start)
	// fmt.Println("s\t", duration)

	// // serialization + write
	// start = time.Now()
	// for c := 0; c < 100000; c++ {
	// 	//byts := encoder.Serialize(testStructs)
	// 	ioutil.WriteFile(fileName, encoder.Serialize(testStructs), os.FileMode(644))
	// }
	// duration = time.Since(start)
	// fmt.Println("s+w\t", duration)

	// // write only
	// byts = encoder.Serialize(testStructs)
	// start = time.Now()
	// for c := 0; c < 100000; c++ {
	// 	ioutil.WriteFile(fileName, byts, os.FileMode(644))
	// }
	// duration = time.Since(start)
	// fmt.Println("w\t", duration)

	// // read only
	// start := time.Now()
	// for c := 0; c < 1000; c++ {
	// 	_, _ = ioutil.ReadFile(fileName)
	// }
	// duration := time.Since(start)
	// fmt.Println("read\t\t", duration)

}
