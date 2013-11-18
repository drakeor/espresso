package rand

const (
	mtSize   = 624
	mtPeriod = 397
	mtDiff   = mtSize - mtPeriod
)

type MT struct {
	Data  []uint32
	Index uint32
}

var MATRIX [2]uint32
var y, i uint32

// MACROS

func m32(val uint32) uint32 {
	return (0x80000000 & val)
}

func l31(val uint32) uint32 {
	return (0x7FFFFFFF & val)
}

func isOdd(val uint32) uint32 {
	return (val & 1)
}

func unroll(M *MT, i *uint32, y *uint32) {
	*y = m32(M.Data[*i]) | l31(M.Data[(*i)+1])
	M.Data[*i] = M.Data[*i-mtDiff] ^ (*y >> 1) ^ MATRIX[isOdd(*y)]
	(*i)++
}

func CreateRandom() *MT {
	mt := &MT{}
	mt.Data = make([]uint32, mtSize)
	mt.Index = 0
	return mt
}

func (M *MT) GenerateNumbers() {
	MATRIX[0] = 0
	MATRIX[1] = 0x9908b0df

	for i = 0; i < mtDiff-1; i++ {
		y = m32(M.Data[i]) | l31(M.Data[i+1])
		M.Data[i] = M.Data[i+mtPeriod] ^ (y >> 1) ^ MATRIX[isOdd(y)]

		i++

		y = m32(M.Data[i]) | l31(M.Data[i+1])
		M.Data[i] = M.Data[i+mtPeriod] ^ (y >> 1) ^ MATRIX[isOdd(y)]
	}

	for i = mtDiff; i < (mtSize - 1); {
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
		unroll(M, &i, &y)
	}

	y = m32(M.Data[mtSize-1]) | l31(M.Data[mtSize-1])
	M.Data[mtSize-1] = M.Data[mtPeriod-1] ^ (y >> 1) ^ MATRIX[isOdd(y)]
}

func (M *MT) Seed(value uint32) {
	M.Data[0] = value
	M.Index = 0

	for i = 1; i < mtSize; i++ {
		M.Data[i] = 0x6c078965*(M.Data[i-1]^M.Data[i-1]>>30) + i
	}
}

func (M *MT) RandU32() uint32 {
	if M.Index == 0 {
		M.GenerateNumbers()
	}

	y = M.Data[M.Index]

	y ^= y >> 11
	y ^= y << 7 & 0x9d2c5680
	y ^= y << 15 & 0xefc60000
	y ^= y >> 18

	M.Index++
	if M.Index == mtSize {
		M.Index = 0
	}

	return y
}
