//+build cxstrat

package model

import (
	. "github.com/SkycoinProject/cx/cx"
	//. "github.com/SkycoinProject/cx/cx/base"
)

const (
	OP_CXDATUM_BEGIN = iota + 65536

	//utility operations
	OP_CXDATUM_STR2BYTES
	OP_CXDATUM_BYTES2STR
	OP_CXDATUM_SUMSHA256
	OP_CXDATUM_RDADDRESS
	OP_CXDATUM_BTADDRESS

	//blockchain operations
	OP_CXDATUM_GETBLOCK
	OP_CXDATUM_CHAINLEN
	OP_CXDATUM_TRANSACT
	OP_CXDATUM_MKACCOUNT
	OP_CXDATUM_LAUNCH

	//boltdb operations
	OP_CXDATUM_NEWBUCKET
	OP_CXDATUM_FETCH
	OP_CXDATUM_STORE

	OP_CXDATUM_STALL
	OP_CXDATUM_EXPOSE
	OP_CXDATUM_LAUNCHAPI

	END_OF_CXDATUM_OPS
)

func init() {
	Op(OP_CXDATUM_STR2BYTES, "cxdatum.str2bytes", str2Bytes, In(ASTR), Out(Slice(TYPE_UI8)))
	Op(OP_CXDATUM_BYTES2STR, "cxdatum.bytes2str", bytes2Str, In(Slice(TYPE_UI8)), Out(ASTR))
	Op(OP_CXDATUM_SUMSHA256, "cxdatum.sumsha256", sumSha256, In(Slice(TYPE_UI8)), Out(Array(TYPE_UI8, 32)))
	Op(OP_CXDATUM_RDADDRESS, "cxdatum.rdaddress", rdAddress, In(Array(TYPE_UI8, 25)), Out(ASTR))
	Op(OP_CXDATUM_BTADDRESS, "cxdatum.btaddress", btAddress, In(Slice(TYPE_UI8)), Out(Array(TYPE_UI8, 25)))

	Op(OP_CXDATUM_GETBLOCK, "cxdatum.getblock", getBlock, In(AI32), Out(Slice(TYPE_UI8)))
	Op(OP_CXDATUM_CHAINLEN, "cxdatum.chainlen", chainLen, nil, Out(AI32))
	Op(OP_CXDATUM_TRANSACT, "cxdatum.transact", transact, In(Slice(TYPE_UI8)), Out(ABOOL))
	Op(OP_CXDATUM_LAUNCH, "cxdatum.launch", launchCXStrat, nil, Out(ABOOL))
	Op(OP_CXDATUM_MKACCOUNT, "cxdatum.mkaccount", makeAccount, nil, Out(ABOOL))

	Op(OP_CXDATUM_NEWBUCKET, "cxdatum.bucket", newBucket, In(ASTR), Out(ABOOL))
	Op(OP_CXDATUM_FETCH, "cxdatum.fetch", fetch, In(ASTR, Slice(TYPE_UI8)), Out(Slice(TYPE_UI8)))
	Op(OP_CXDATUM_STORE, "cxdatum.store", store, In(ASTR, Slice(TYPE_UI8), Slice(TYPE_UI8)), Out(ABOOL))

	Op(OP_CXDATUM_STALL, "cxdatum.stall", stall, nil, Out(Slice(TYPE_UI8)))
	Op(OP_CXDATUM_EXPOSE, "cxdatum.expose", expose, In(Slice(TYPE_UI8)), nil)
	Op(OP_CXDATUM_LAUNCHAPI, "cxdatum.launchapi", launchApi, nil, Out(ABOOL))
}
