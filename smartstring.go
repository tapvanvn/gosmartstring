package gosmartstring

import "github.com/tapvanvn/gotokenize"

var (
	//MARK: smartstring
	TokenSSLOperator = 1
	TokenSSLString   = 2

	TokenSSLParenthese = 100
	TokenSSLBlock      = 101
	TokenSSLSquare     = 102

	TokenSSLCommand = 1000

	TokenSSLSmarstring   = 1001
	TokenSSLNormalstring = 1002

	//MARK:
	TokenSSRegistryIgnore = 10 //dont care result
	TokenSSRegistry       = 11 //link to registry
	TokenSSRegistryGlobal = 12 //set result registry address to global

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

var SSInstructionTokenMove int = 0

func SSInsructionMove(delta int) {

	TokenSSInstructionDo += delta
	TokenSSInstructionLink += delta
	TokenSSInstructionPack += delta
	TokenSSInstructionExport += delta
	TokenSSInstructionIf += delta
	TokenSSInstructionCase += delta
	TokenSSInstructionEach += delta
	TokenSSInstructionCount += delta

}

var SSLGlobalNested = []int{
	TokenSSLSmarstring,
}
var SSLIgnores = []int{}

var SSLPatterns = []gotokenize.Pattern{
	//pattern attribute "key"="value"
	{
		Type: TokenSSLCommand,
		Struct: []gotokenize.PatternToken{
			{Type: 0},
			{Type: TokenSSLParenthese, CanNested: true},
		},
		IsRemoveGlobalIgnore: true,
	},
}
