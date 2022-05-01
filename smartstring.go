package gosmartstring

import "github.com/tapvanvn/gotokenize/v2"

var (
	//MARK: smartstring
	TokenSSLOperator = 1
	TokenSSLString   = 2

	TokenSSLParenthese = 10
	TokenSSLBlock      = 11
	TokenSSLSquare     = 12

	TokenSSLCommand      = 100
	TokenSSLInstruction  = 101
	TokenSSLSmartstring  = 102
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
	TokenSSLSmartstring += delta
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
	TokenSSLSmartstring,
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

func SSNaming(tokenType int) string {
	switch tokenType {
	case TokenSSLOperator:
		return "ss_operator"
	case TokenSSLString:
		return "ss_string"
	case TokenSSLParenthese:
		return "ss_parenthese"
	case TokenSSLBlock:
		return "ss_block"
	case TokenSSLSquare:
		return "ss_square"
	case TokenSSLCommand:
		return "ss_command"
	case TokenSSLInstruction:
		return "ss_instruction"
	case TokenSSLSmartstring:
		return "ss_smartstring"
	case TokenSSLNormalstring:
		return "ss_normalstring"
	case TokenSSRegistryIgnore:
		return "ss_registry_ignore"
	case TokenSSRegistry:
		return "ss_registry"
	case TokenSSRegistryGlobal:
		return "ss_registry_global"
	case TokenSSInstructionDo:
		return "ss_do"
	case TokenSSInstructionLink:
		return "ss_link"
	case TokenSSInstructionPack:
		return "ss_pack"
	case TokenSSInstructionExport:
		return "ss_export"
	case TokenSSInstructionIf:
		return "ss_if"
	case TokenSSInstructionCase:
		return "ss_case"
	case TokenSSInstructionEach:
		return "ss_each"
	case TokenSSInstructionCount:
		return "ss_count"
	default:
		return "unknown"
	}
}

var EmptyParams []IObject = []IObject{}
