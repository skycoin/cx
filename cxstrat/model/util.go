//+build cxstrat

package model

import (
	"regexp"

	. "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/skycoin/src/cipher"
)

func str2Bytes(prgm *CXProgram) {
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	input1, output1 := expr.Inputs[0], expr.Outputs[0]

	byts := []byte(ReadStr(fp, input1))

	outputSlicePointer := GetFinalOffset(fp, output1)
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(byts)), 1))
	copy(GetSliceData(outputSliceOffset, 1), byts)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func bytes2Str(prgm *CXProgram) {
	/* works special now */
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	byts := GetSliceData(inputSliceOffset, 1)

	nbyts := []byte{}
	for i := 0; i < len(byts) && byts[i] != 0; i++ {
		nbyts = append(nbyts, byts[i])
	}

	str := string(nbyts)
	//fmt.Println(str)
	WriteString(fp, str, expr.Outputs[0])
}

func sumSha256(prgm *CXProgram) {
	//expects ui8 slice, returns [32]ui8 hash.
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	byts := GetSliceData(inputSliceOffset, 1)
	hsh := cipher.SumSHA256(byts)

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), hsh[:])
}

func rdAddress(prgm *CXProgram) {
	//expects ui8 slice, returns [32]ui8 hash.
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	abyts := ReadMemory(GetFinalOffset(fp, expr.Inputs[0]), expr.Inputs[0])

	a, _ := cipher.AddressFromBytes(abyts)
	astr := a.String()

	WriteString(fp, astr, expr.Outputs[0])
}

func btAddress(prgm *CXProgram) {
	//expects ui8 slice, returns [32]ui8 hash.
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	byts := GetSliceData(inputSliceOffset, 1)

	abyts := string(byts)
	a, _ := cipher.DecodeBase58Address(abyts)
	astr := a.Bytes()

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), astr)
}

func makeHashIdx(t string) []int {
	var rv []int
	// := []byte(t_)
	for i := 0; i < len(t); i++ {
		if t[i] == '#' {
			i++
			/* store idx */
			rv = append(rv, i)

			if i == len(t) {
				return rv
			}
		}
	}

	return rv
}

func makeHashs(t string) []string {
	var rv []string
	// := []byte(t_)
	regt := regexp.MustCompile("[0-9A-Za-z_]")
	for i := 0; i < len(t); i++ {
		if t[i] == '#' {
			i++
			var tstring []byte
			for regt.Match([]byte{byte(t[i])}) && i < len(t) {
				tstring = append(tstring, byte(t[i]))
				i++
			}

			//fmt.Println("I am a happy hashtag: " + string(tstring))
			rv = append(rv, string(tstring))
			if i == len(t) {
				return rv
			}
		}
	}

	return rv
}

func makeTags(t string) []cipher.Address {
	var rv []cipher.Address
	// := []byte(t_)
	regx := regexp.MustCompile("[0-9A-Za-z_]")
	for i := 0; i < len(t); i++ {
		if t[i] == '@' {
			i++
			var tstring []byte
			for regx.Match([]byte{byte(t[i])}) && i < len(t) {
				tstring = append(tstring, byte(t[i]))
				i++
			}
			if i == len(t) {
				return rv
			}

			taddr_, err := cipher.DecodeBase58Address(string(tstring))
			//fmt.Println("I am a happy wallet: " + taddr_.String())
			if err != nil {
				//cache.err = ERRORCODE_TWEETINVALID
				return rv
			}
			rv = append(rv, taddr_)
		}
	}

	return rv
}

func waitForAPIReturn() []byte {
	<-donereq
	return <-req
}

/**
 * 1. Program calls cxtweet.stall(), freeing control to API.
 * 2. API receives a call to an endpoint (functionXXX).
 * 3. functionXXX fills channel with received object.
 * 4. waitUntilAPIReturn() uses select to block until this.
 * 5. ^^^ returns bytes array to program.
 * 6. Program checks byte array to determine kind of API handler to use.
 * 7. Program serves.
 * 8. Program calls expose() to expose value and fill a channel.
 * 9. waitUntilPrgmReturn() was called by endpoint.
 * A. ^ uses select to block until channel is filled (which it is now).
 * B. Passes channel buffer to functionXXX.
 * C. functionXXX can finally return, deserializes channel buffer to return object.
 * D. GOTO step 1.
 */
func stall(prgm *CXProgram) {
	/* stalls the program until value is returned. */
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	byts := waitForAPIReturn()

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(byts)), 1))
	copy(GetSliceData(outputSliceOffset, 1), byts)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func expose(prgm *CXProgram) {
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	byts := GetSliceData(inputSliceOffset, 1)

	res <- byts
	doneres <- true
}
