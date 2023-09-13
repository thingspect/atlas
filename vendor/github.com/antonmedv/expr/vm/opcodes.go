package vm

type Opcode byte

const (
	OpInvalid Opcode = iota
	OpPush
	OpInt
	OpPop
	OpStore
	OpLoadVar
	OpLoadConst
	OpLoadField
	OpLoadFast
	OpLoadMethod
	OpLoadFunc
	OpLoadEnv
	OpFetch
	OpFetchField
	OpMethod
	OpTrue
	OpFalse
	OpNil
	OpNegate
	OpNot
	OpEqual
	OpEqualInt
	OpEqualString
	OpJump
	OpJumpIfTrue
	OpJumpIfFalse
	OpJumpIfNil
	OpJumpIfNotNil
	OpJumpIfEnd
	OpJumpBackward
	OpIn
	OpLess
	OpMore
	OpLessOrEqual
	OpMoreOrEqual
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
	OpExponent
	OpRange
	OpMatches
	OpMatchesConst
	OpContains
	OpStartsWith
	OpEndsWith
	OpSlice
	OpCall
	OpCall0
	OpCall1
	OpCall2
	OpCall3
	OpCallN
	OpCallFast
	OpCallTyped
	OpCallBuiltin1
	OpArray
	OpMap
	OpLen
	OpCast
	OpDeref
	OpIncrementIndex
	OpDecrementIndex
	OpIncrementCount
	OpGetIndex
	OpSetIndex
	OpGetCount
	OpGetLen
	OpGetGroupBy
	OpGetAcc
	OpPointer
	OpThrow
	OpGroupBy
	OpSetAcc
	OpBegin
	OpEnd // This opcode must be at the end of this list.
)
