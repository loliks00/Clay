package main

import (
	"fmt"
	"os"
)

type TokenType int8

const (
	TokenType_Plus TokenType = iota
	TokenType_Minus
	TokenType_Divide
	TokenType_Multiply
	TokenType_Modulo

	TokenType_Inc
	TokenType_Dec

	TokenType_Grth
	TokenType_GrthEq
	TokenType_Less
	TokenType_LessEq
	TokenType_Eq
	TokenType_NotEq

	TokenType_Not //!
	TokenType_And // |
	TokenType_Or  // &

	TokenType_BinAnd // &&
	TokenType_BinOr  // ||
	TokenType_BinXor //~~
	TokenType_BinNot //!!
	TokenType_BinShl
	TokenType_BinShr

	TokenType_Comma
	TokenType_Dot
	TokenType_DotDotDot
	TokenType_Colon
	TokenType_Assign
	TokenType_AssingInfer //:=
	TokenType_ColonColon

	TokenType_Deref     //^
	TokenType_At        // @
	TokenType_Arrow     // ->
	TokenType_Dolarsign // $

	TokenType_Keyword_If
	TokenType_Keyword_Else
	TokenType_Keyword_For
	TokenType_Keyword_Switch
	TokenType_Keyword_Fallthrough
	TokenType_Keyword_Break
	TokenType_Keyword_Continue
	TokenType_Keyword_Import
	TokenType_Keyword_Defer
	TokenType_Keyword_Default
	TokenType_Keyword_Range
	TokenType_Keyword_Enum
	TokenType_Keyword_Struct
	TokenType_Keyword_S8
	TokenType_Keyword_S16
	TokenType_Keyword_S32
	TokenType_Keyword_S64
	TokenType_Keyword_U8
	TokenType_Keyword_U16
	TokenType_Keyword_U32
	TokenType_Keyword_U64
	TokenType_Keyword_F32
	TokenType_Keyword_F64
	TokenType_Keyword_Bool
	TokenType_Keyword_String
	TokenType_Keyword_CString
	TokenType_Keyword_Map
	TokenType_Keyword_Typeid
	TokenType_Keyword_Rune
	TokenType_Keyword_Byte
	TokenType_Keyword_True
	TokenType_Keyword_False

	TokenType_Directive // for example #pure
	TokenType_Lable     // for example lable#
	TokenType_Symbol
	TokenType_IntLit
	TokenType_IntLit_Hex
	TokenType_IntLit_Bin
	TokenType_IntLit_Oct
	TokenType_FloatLit
	TokenType_StringLit
	TokenType_StringLongLit

	TokenType_LParen
	TokenType_RParen
	TokenType_LCurly
	TokenType_RCurly
	TokenType_LBrack
	TokenType_RBrack

	TokenType_EOF
	TokenType_NewLine
)

type Token struct {
	Pos  uint
	Len  uint16
	Type TokenType
}

type Parser struct {
	Tokens []Token
	Pos    int
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}
func isChar(r rune) bool {
	return (r >= 'a' && r <= 'z') && (r >= 'A' && r <= 'Z') && r == '_'
}

func (p *Parser) Lex(src []rune) {
	line := 0
	for i := uint(0); i < uint(len(src)); i++ {
		switch src[i] {
		case ' ':
			continue
		case '\r':
			continue
		case '\n':
			line++
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_NewLine})
		case 0:
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_EOF})
			return
		case '{':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_LCurly})
			continue
		case '}':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_RCurly})
			continue
		case '[':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_LBrack})
			continue
		case ']':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_RBrack})
			continue
		case '(':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_LParen})
			continue
		case ')':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_RParen})
			continue
		case '+':
			if i+1 < uint(len(src)) && src[i+1] == '+' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Inc})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Plus})
			}
			continue
		case '-':
			if i+1 < uint(len(src)) && src[i+1] == '-' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Dec})
			} else if i+1 < uint(len(src)) && src[i+1] == '>' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Arrow})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Minus})
			}
			continue
		case '*':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Divide})
			continue
		case '/':
			if i+1 < uint(len(src)) && src[i+1] == '/' {
				i = i + 2
				for !(src[i] == 0) && !(src[i] == '\n') {
					i++
				}
			} else if i+1 < uint(len(src)) && src[i+1] == '*' {
				i = i + 2
				for !(src[i] == 0 || src[i+1] == 0) && !(src[i] == '*' && src[i+1] == '/') {
					i++
				}
				i++
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Divide})
			}
			continue
		case '%':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Modulo})
			continue
		case '!':
			if i+1 < uint(len(src)) && src[i+1] == '!' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_BinNot})
			} else if i+1 < uint(len(src)) && src[i+1] == '=' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_NotEq})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Not})
			}
			continue
		case '<':
			if i+1 < uint(len(src)) && src[i+1] == '=' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_LessEq})
			} else if i+1 < uint(len(src)) && src[i+1] == '<' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_BinShl})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Less})
			}
			continue
		case '>':
			if i+1 < uint(len(src)) && src[i+1] == '=' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_GrthEq})
			} else if i+1 < uint(len(src)) && src[i+1] == '>' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_BinShr})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Grth})
			}
			continue
		case ':':
			if i+1 < uint(len(src)) && src[i+1] == '=' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_AssingInfer})
			} else if i+1 < uint(len(src)) && src[i+1] == ':' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_ColonColon})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Colon})
			}
			continue
		case '=':
			if i+1 < uint(len(src)) && src[i+1] == '=' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Eq})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Assign})
			}
			continue
		case ',':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Comma})
			continue
		case '.':
			if i+2 < uint(len(src)) && src[i+1] == '.' && src[i+2] == '.' {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_DotDotDot})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Dot})
			}
			continue
		case '&':
			if i+1 < uint(len(src)) && src[i+1] == '&' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_BinAnd})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_And})
			}
			continue
		case '^':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Deref})
			continue
		case '@':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_At})
			continue
		case '~':
			if i+1 < uint(len(src)) && src[i+1] == '~' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_BinXor})
				continue
			}
		case '|':
			if i+1 < uint(len(src)) && src[i+1] == '|' {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_BinOr})
			} else {
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Or})
			}
			continue
		case '$':
			p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Dolarsign})
			continue

		case 'i':
			if i+1 < uint(len(src)) && src[i+1] == 'f' && (!isChar(src[i+2]) && !isNumber(src[i+2])) {
				i++
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_If})
				continue
			}
			if i+5 < uint(len(src)) && src[i+1] == 'm' && src[i+2] == 'p' && src[i+3] == 'o' && src[i+4] == 'r' && src[i+5] == 't' && (!isChar(src[i+6]) && !isNumber(src[i+6])) {
				i = i + 5
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Import})
				continue
			}
		case 'e':
			if i+3 < uint(len(src)) && src[i+1] == 'l' && src[i+2] == 's' && src[i+3] == 'e' && (!isChar(src[i+4]) && !isNumber(src[i+4])) {
				i = i + 3
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Else})
				continue
			} else if i+3 < uint(len(src)) && src[i+1] == 'n' && src[i+2] == 'u' && src[i+3] == 'm' && (!isChar(src[i+4]) && !isNumber(src[i+4])) {
				i = i + 3
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Enum})
				continue
			}

		case 'f':
			if i+2 < uint(len(src)) && src[i+1] == 'o' && src[i+2] == 'r' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_For})
				continue
			} else if i+10 < uint(len(src)) && src[i+1] == 'a' && src[i+2] == 'l' && src[i+3] == 'l' && src[i+4] == 't' && src[i+5] == 'h' && src[i+6] == 'r' && src[i+7] == 'o' && src[i+8] == 'u' && src[i+9] == 'g' && src[i+10] == 'h' && (!isChar(src[i+11]) && !isNumber(src[i+11])) {
				i = i + 10
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Fallthrough})
				continue
			} else if i+4 < uint(len(src)) && src[i+1] == 'a' && src[i+2] == 'l' && src[i+3] == 's' && src[i+4] == 'e' && (!isChar(src[i+5]) && !isNumber(src[i+5])) {
				i = i + 4
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_False})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '3' && src[i+2] == '2' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_F32})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '6' && src[i+2] == '4' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_F64})
				continue
			}

		case 's':
			if i+5 < uint(len(src)) && src[i+1] == 'w' && src[i+2] == 'i' && src[i+3] == 't' && src[i+4] == 'c' && src[i+5] == 'h' && (!isChar(src[i+6]) && !isNumber(src[i+6])) {
				i = i + 5
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Switch})
				continue
			} else if i+5 < uint(len(src)) && src[i+1] == 't' && src[i+2] == 'r' && src[i+3] == 'u' && src[i+4] == 'c' && src[i+5] == 't' && (!isChar(src[i+6]) && !isNumber(src[i+6])) {
				i = i + 5
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Struct})
				continue
			} else if i+5 < uint(len(src)) && src[i+1] == 't' && src[i+2] == 'r' && src[i+3] == 'i' && src[i+4] == 'n' && src[i+5] == 'g' && (!isChar(src[i+6]) && !isNumber(src[i+6])) {
				i = i + 5
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_String})
				continue
			} else if i+1 < uint(len(src)) && src[i+1] == '8' && (!isChar(src[i+2]) && !isNumber(src[i+2])) {
				i = i + 1
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_S8})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '1' && src[i+2] == '6' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_S16})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '3' && src[i+2] == '2' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_S32})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '6' && src[i+2] == '4' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_S64})
				continue
			}

		case 'b':
			if i+4 < uint(len(src)) && src[i+1] == 'r' && src[i+2] == 'e' && src[i+3] == 'a' && src[i+4] == 'k' && (!isChar(src[i+5]) && !isNumber(src[i+5])) {
				i = i + 4
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Break})
				continue
			} else if i+3 < uint(len(src)) && src[i+1] == 'y' && src[i+2] == 't' && src[i+3] == 'e' && (!isChar(src[i+4]) && !isNumber(src[i+4])) {
				i = i + 3
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Byte})
				continue
			} else if i+3 < uint(len(src)) && src[i+1] == 'o' && src[i+2] == 'o' && src[i+3] == 'l' && (!isChar(src[i+4]) && !isNumber(src[i+4])) {
				i = i + 3
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Bool})
				continue
			}
		case 'c':
			if i+7 < uint(len(src)) && src[i+1] == 'o' && src[i+2] == 'n' && src[i+3] == 't' && src[i+4] == 'i' && src[i+5] == 'n' && src[i+6] == 'u' && src[i+7] == 'e' && (!isChar(src[i+8]) && !isNumber(src[i+8])) {
				i = i + 7
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Continue})
				continue
			} else if i+6 < uint(len(src)) && src[i+1] == 's' && src[i+2] == 't' && src[i+3] == 'r' && src[i+4] == 'i' && src[i+5] == 'n' && src[i+6] == 'g' && (!isChar(src[i+7]) && !isNumber(src[i+7])) {
				i = i + 6
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_CString})
				continue
			}
		case 'd':
			if i+4 < uint(len(src)) && src[i+1] == 'e' && src[i+2] == 'f' && src[i+3] == 'e' && src[i+4] == 'r' && (!isChar(src[i+5]) && !isNumber(src[i+5])) {
				i = i + 4
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Defer})
				continue
			} else if i+6 < uint(len(src)) && src[i+1] == 'e' && src[i+2] == 'f' && src[i+3] == 'a' && src[i+4] == 'u' && src[i+5] == 'l' && src[i+6] == 't' && (!isChar(src[i+7]) && !isNumber(src[i+7])) {
				i = i + 6
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Default})
				continue
			}
		case 'r':
			if i+4 < uint(len(src)) && src[i+1] == 'a' && src[i+2] == 'n' && src[i+3] == 'g' && src[i+4] == 'e' && (!isChar(src[i+5]) && !isNumber(src[i+5])) {
				i = i + 4
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Range})
				continue
			} else if i+3 < uint(len(src)) && src[i+1] == 'u' && src[i+2] == 'n' && src[i+3] == 'e' && (!isChar(src[i+4]) && !isNumber(src[i+4])) {
				i = i + 3
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Rune})
				continue
			}
		case 't':
			if i+3 < uint(len(src)) && src[i+1] == 'r' && src[i+2] == 'u' && src[i+3] == 'e' && (!isChar(src[i+4]) && !isNumber(src[i+4])) {
				i = i + 3
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_True})
				continue
			} else if i+5 < uint(len(src)) && src[i+1] == 'y' && src[i+2] == 'p' && src[i+3] == 'e' && src[i+4] == 'i' && src[i+5] == 'd' && (!isChar(src[i+5]) && !isNumber(src[i+5])) {
				i = i + 5
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Typeid})
				continue
			}
		case 'm':
			if i+2 < uint(len(src)) && src[i+1] == 'a' && src[i+2] == 'p' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_Map})
				continue
			}
		case 'u':
			if i+2 < uint(len(src)) && src[i+1] == '8' && (!isChar(src[i+2]) && !isNumber(src[i+2])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_U8})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '6' && src[i+2] == '4' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_U16})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '6' && src[i+2] == '4' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_U32})
				continue
			} else if i+2 < uint(len(src)) && src[i+1] == '6' && src[i+2] == '4' && (!isChar(src[i+3]) && !isNumber(src[i+3])) {
				i = i + 2
				p.Tokens = append(p.Tokens, Token{Pos: i, Len: 0, Type: TokenType_Keyword_U64})
				continue
			}
		}
		if src[i] == '#' {
			pos := i
			leng := 0
			for isChar(src[i]) && i+2 > uint(len(src)) {
				i++
				leng++
			}
			p.Tokens = append(p.Tokens, Token{Pos: pos, Len: uint16(leng), Type: TokenType_Directive})
			continue

		} else if src[i] >= '0' && src[i] <= '9' {
			if src[i+1] == '0' {
				pos := i
				i++
				switch src[i+1] {
				case 'o':
					i++
					for (isNumber(src[i]) || src[i] == '\'') && uint(len(src)) > i {
						i++
					}
				case 'x':
				case 'b':
				}
			}
			isFloat := false

		} else if (src[i] >= 'a' && src[i] <= 'z') && (src[i] >= 'A' && src[i] <= 'Z') && src[i] == '_' {

		} else if src[i] == '"' {

		} else if src[i] == '\'' {

		} else {
			fmt.Printf("Unknown Token '%c'on line : %d", src[i], line)
			os.Exit(1)
		}

	}
}

func main() {}

/*



	context :: struct {
		base : rawptr // ptr used for relative pointers
		allocator : <- allocator
		logger : logger

	}



	context context is inmutable,
	context is scoped and passed to functions



	func :: (){
		x~ = 123 //base is some curent ptr
		func2() // base is changed
		println("%d",x~) //base is returned to normal
	}

	func2 :: (){
		context.base = somerawptr
	}


	~~ <- not

*/
