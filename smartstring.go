package gosmartstring

import "github.com/tapvanvn/gotokenize"

var (
	//MARK: smartstring
	TokenSSLOperator = 1
	TokenSSLString   = 2

	TokenSSLParenthese = 10
	TokenSSLBlock      = 11
	TokenSSLSquare     = 12

	TokenSSLCommand      = 100
	TokenSSLInstruction  = 101
	TokenSSLSmarstring   = 102
	TokenSSLNormalstring = 103

	//MARK:
	TokenSSRegistryIgnore = 200 //dont care result
	TokenSSRegistry       = 201 //link to registry
	TokenSSRegistryGlobal = 202 //set result registry address to global

	TokenSSInstructionDo       = 300 //command to do
	TokenSSInstructionLink     = 301 //link last instruction to be input of next instruction
	TokenSSInstructionRemember = 302 //set remember flag to true
	TokenSSInstructionPack     = 303 //each children is an instruction
	TokenSSInstructionExport   = 304 //just export string
	TokenSSInstructionIf       = 305 //if statement
	TokenSSInstructionCase     = 306 //check in cases
	TokenSSInstructionEach     = 307 //loop for each .. in .. and do
	TokenSSInstructionCount    = 308 //count to and do
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

	//MARK: smartstring
	TokenSSLOperator += delta
	TokenSSLString += delta

	TokenSSLParenthese += delta
	TokenSSLBlock += delta
	TokenSSLSquare += delta

	TokenSSLCommand += delta
	TokenSSLInstruction += delta
	TokenSSLSmarstring += delta
	TokenSSLNormalstring += delta

	//MARK:
	TokenSSRegistryIgnore += delta
	TokenSSRegistry += delta
	TokenSSRegistryGlobal += delta

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
