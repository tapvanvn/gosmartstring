package gosmartstring

import "github.com/tapvanvn/gotokenize/v2"

var (
	//MARK: smartstring
	TokenSSLOperator = 1
	TokenSSLString   = 2
	TokenSSLWord     = 3

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

	TokenSSInstructionDo     = 300 //command to do
	TokenSSInstructionLink   = 301 //link last instruction to be input of next instruction
	TokenSSInstructionReload = 302 //check if can reload last returned object. error if last returned object is nil
	TokenSSInstructionPack   = 303 //each children is an instruction
	TokenSSInstructionExport = 304 //just export string
	TokenSSInstructionIf     = 305 //if statement
	TokenSSInstructionCase   = 306 //check in cases
	TokenSSInstructionEach   = 307 //loop for each .. in .. and do
	TokenSSInstructionCount  = 308 //count to and do
)

var SSLAllTokens = []*int{
	&TokenSSLOperator,
	&TokenSSLString,
	&TokenSSLWord,
	&TokenSSLParenthese,
	&TokenSSLBlock,
	&TokenSSLSquare,
	&TokenSSLCommand,
	&TokenSSLInstruction,
	&TokenSSLSmartstring,
	&TokenSSLNormalstring,
	&TokenSSRegistryIgnore,
	&TokenSSRegistry,
	&TokenSSRegistryGlobal,
	&TokenSSInstructionDo,
	&TokenSSInstructionLink,
	&TokenSSInstructionReload,
	&TokenSSInstructionPack,
	&TokenSSInstructionExport,
	&TokenSSInstructionIf,
	&TokenSSInstructionCase,
	&TokenSSInstructionEach,
	&TokenSSInstructionCount,
}

var SSInstructionTokenMove int = 0

func SSInsructionMove(delta int) {

	for _, token := range SSLAllTokens {
		*token += delta
	}
}

var SSLIgnores = []int{}

func getSSLGlobalNested() []int {
	return []int{
		TokenSSLSmartstring,
		TokenSSLParenthese,
	}
}
func buildSSLPatterns() []gotokenize.Pattern {
	return []gotokenize.Pattern{

		{
			Type: TokenSSLCommand,
			Struct: []gotokenize.PatternToken{
				{Type: TokenSSLWord},
				{Type: TokenSSLParenthese, CanNested: true},
			},
			IsRemoveGlobalIgnore: true,
		},
	}
}

func SSNaming(tokenType int) string {
	switch tokenType {
	case TokenSSLOperator:
		return "ss_operator"
	case TokenSSLString:
		return "ss_string"
	case TokenSSLWord:
		return "ss_word"
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
	case TokenSSInstructionReload:
		return "ss_reload"
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
