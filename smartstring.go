package gosmartstring

const (
	TokenSSRegistry = 10 //link to registry

	TokenSSInstructionDo     = 100 //command to do
	TokenSSInstructionLink   = 101 //link last instruction to be input of next instruction
	TokenSSInstructionPack   = 102 //each children is an instruction
	TokenSSInstructionExport = 103 //just export string
	TokenSSInstructionIf     = 104 //if statement
	TokenSSInstructionCase   = 105 //check in cases
	TokenSSInstructionEach   = 106 //loop for each .. in .. and do
	TokenSSInstructionCount  = 107 //count to and do
)

var AllSSInstructions = []int{
	TokenSSInstructionDo,
	TokenSSInstructionLink,
	TokenSSInstructionPack,
	TokenSSInstructionExport,
	TokenSSInstructionIf,
	TokenSSInstructionCase,
	TokenSSInstructionEach,
	TokenSSInstructionCount,
}
