package gosmartstring

import "github.com/tapvanvn/gotokenize/v2"

var (
	//MARK: smartstring
	TokenSSLOperator = 1
	TokenSSLString   = 2
	TokenSSLWord     = 3
	TokenSSLBreak    = 4

	TokenSSLParenthese = 10
	TokenSSLBlock      = 11
	TokenSSLSquare     = 12

	TokenSSLCommand      = 100
	TokenSSLSmartstring  = 102
	TokenSSLNormalstring = 103
	TokenSSLPair         = 104

	//MARK:
	TokenSSRegistryIgnore = 200 //dont care result
	TokenSSRegistry       = 201 //link to registry
	TokenSSRegistryGlobal = 202 //set result registry address to global

	TokenSSInstructionDo               = 300 //command to do
	TokenSSInstructionLink             = 301 //link last instruction to be input of next instruction
	TokenSSInstructionReload           = 302 //check if can reload last returned object. error if last returned object is nil
	TokenSSInstructionPack             = 303 //each children is an instruction
	TokenSSInstructionExport           = 304 //just export string
	TokenSSInstructionIf               = 305 //if statement
	TokenSSInstructionCase             = 306 //check in cases
	TokenSSInstructionEach             = 307 //loop for each .. in .. and do
	TokenSSInstructionCount            = 308 //count to and do
	TokenSSInstructionReset            = 309 //reset this
	TokenSSInstructionQuestion         = 310 //question to skip
	TokenSSInstructionNegativeQuestion = 311 //question (negative) to skip
	TokenSSInstructionBuildObject      = 312
)

var SSLAllTokens = []*int{
	&TokenSSLOperator,
	&TokenSSLString,
	&TokenSSLWord,
	&TokenSSLBreak,
	&TokenSSLParenthese,
	&TokenSSLBlock,
	&TokenSSLSquare,
	&TokenSSLCommand,
	&TokenSSLSmartstring,
	&TokenSSLNormalstring,
	&TokenSSLPair,
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
	&TokenSSInstructionReset,
	&TokenSSInstructionBuildObject,
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
		TokenSSLSquare,
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
	case TokenSSLBreak:
		return "ss_break"
	case TokenSSLParenthese:
		return "ss_parenthese"
	case TokenSSLBlock:
		return "ss_block"
	case TokenSSLSquare:
		return "ss_square"
	case TokenSSLCommand:
		return "ss_command"
	case TokenSSLPair:
		return "pair"
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
	case TokenSSInstructionReset:
		return "ss_reset"
	default:
		return "unknown"
	}
}

var EmptyParams []IObject = []IObject{}
