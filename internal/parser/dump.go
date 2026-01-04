package parser

import (
	"fmt"
	"strings"
)

func DumpProgram(p *Program) string {
	var out strings.Builder
	for i, e := range p.Expressions {
		out.WriteString(dumpExpr(e, "", i == len(p.Expressions)-1))
	}
	return out.String()
}

func node(indent string, last bool, label string) (string, string) {
	branch := "├─ "
	next := indent + "│  "
	if last {
		branch = "└─ "
		next = indent + "   "
	}
	return indent + branch + label + "\n", next
}

func dumpExpr(expr Expression, indent string, last bool) string {
	switch n := expr.(type) {
	case *Identifier:
		line, _ := node(indent, last, "Identifier "+n.Name)
		return line
	case *IntegerLiteral:
		line, _ := node(indent, last, fmt.Sprintf("IntegerLiteral %d", n.Value))
		return line
	case *FloatLiteral:
		line, _ := node(indent, last, fmt.Sprintf("FloatLiteral %g", n.Value))
		return line
	case *StringLiteral:
		line, _ := node(indent, last, fmt.Sprintf("StringLiteral %q", n.Value))
		return line
	case *CharacterLiteral:
		line, _ := node(indent, last, fmt.Sprintf("CharacterLiteral '%c'", n.Value))
		return line
	case *BooleanLiteral:
		line, _ := node(indent, last, fmt.Sprintf("BooleanLiteral %t", n.Value))
		return line
	case *UnitLiteral:
		lex := "()"
		line, _ := node(indent, last, fmt.Sprintf("UnitLiteral %s", lex))
		return line
	case *ListExpression:
		line, next := node(indent, last, "ListExpression")
		var out strings.Builder
		out.WriteString(line)
		for i, e := range n.Elements {
			out.WriteString(dumpExpr(e, next, i == len(n.Elements)-1))
		}
		return out.String()
	case *IndexExpression:
		line, next := node(indent, last, "IndexExpression")
		var out strings.Builder
		out.WriteString(line)
		tLine, tNext := node(next, false, "Target")
		out.WriteString(tLine)
		out.WriteString(dumpExpr(n.Target, tNext, true))
		iLine, iNext := node(next, true, "Index")
		out.WriteString(iLine)
		out.WriteString(dumpExpr(n.Index, iNext, true))
		return out.String()
	case *PrefixExpression:
		line, next := node(indent, last, "PrefixExpression "+n.Operator.Lexeme)
		return line + dumpExpr(n.Right, next, true)
	case *InfixExpression:
		line, next := node(indent, last, "InfixExpression "+n.Operator.Lexeme)
		return line +
			dumpExpr(n.Left, next, false) +
			dumpExpr(n.Right, next, true)
	case *CallExpression:
		line, next := node(indent, last, "CallExpression")
		var out strings.Builder
		out.WriteString(line)
		cLine, cNext := node(next, len(n.Arguments) == 0, "Callee")
		out.WriteString(cLine)
		out.WriteString(dumpExpr(n.Callee, cNext, true))
		if len(n.Arguments) > 0 {
			aLine, aNext := node(next, true, "Arguments")
			out.WriteString(aLine)
			for i, arg := range n.Arguments {
				out.WriteString(dumpExpr(arg, aNext, i == len(n.Arguments)-1))
			}
		}
		return out.String()
	case *VariableDeclarationExpression:
		line, next := node(indent, last,
			fmt.Sprintf("VariableDeclaration name=%s rec=%t", n.Name.Lexeme, n.Recursive))
		var out strings.Builder
		out.WriteString(line)
		if n.Type != nil {
			tLine, tNext := node(next, false, "Type")
			out.WriteString(tLine)
			out.WriteString(dumpType(n.Type, tNext, true))
		}
		vLine, vNext := node(next, true, "Value")
		out.WriteString(vLine)
		out.WriteString(dumpExpr(n.Value, vNext, true))
		return out.String()
	case *FunctionDeclarationExpression:
		line, next := node(indent, last,
			fmt.Sprintf("FunctionDeclaration name=%s rec=%t", n.Name.Lexeme, n.Recursive))
		var out strings.Builder
		out.WriteString(line)
		if n.Signature != nil {
			tLine, tNext := node(next, true, "TypeSignature")
			out.WriteString(tLine)
			out.WriteString(dumpType(n.Signature, tNext, true))
		}
		fLine, fNext := node(next, true, "Function")
		out.WriteString(fLine)
		out.WriteString(dumpExpr(n.Function, fNext, true))
		return out.String()
	case *FunctionLiteralExpression:
		line, next := node(indent, last, "FunctionLiteral")
		var out strings.Builder
		out.WriteString(line)
		if len(n.Parameters) > 0 {
			pLine, pNext := node(next, false, "Parameters")
			out.WriteString(pLine)
			for i, p := range n.Parameters {
				pl, pi := node(pNext, i == len(n.Parameters)-1, "Param "+p.Name.Lexeme)
				out.WriteString(pl)
				if p.Type != nil {
					out.WriteString(dumpType(p.Type, pi, true))
				}
			}
		}
		bLine, bNext := node(next, true, "Body")
		out.WriteString(bLine)
		out.WriteString(dumpExpr(n.Body, bNext, true))
		return out.String()
	case *BlockExpression:
		line, next := node(indent, last, "BlockExpression")
		var out strings.Builder
		out.WriteString(line)
		for i, e := range n.Expressions {
			out.WriteString(dumpExpr(e, next, i == len(n.Expressions)-1))
		}
		return out.String()
	case *IfExpression:
		line, next := node(indent, last, "IfExpression")
		var out strings.Builder
		out.WriteString(line)
		cLine, cNext := node(next, false, "Condition")
		out.WriteString(cLine)
		out.WriteString(dumpExpr(n.Condition, cNext, true))
		tLine, tNext := node(next, n.Else == nil, "Then")
		out.WriteString(tLine)
		out.WriteString(dumpExpr(n.Then, tNext, true))
		if n.Else != nil {
			eLine, eNext := node(next, true, "Else")
			out.WriteString(eLine)
			out.WriteString(dumpExpr(n.Else, eNext, true))
		}
		return out.String()
	case *ImportExpression:
		line, _ := node(indent, last, "Import "+n.Module)
		return line
	case *SliceExpression:
		line, next := node(indent, last, "SliceExpression")
		var out strings.Builder
		out.WriteString(line)
		tLine, tNext := node(next, false, "Target")
		out.WriteString(tLine)
		out.WriteString(dumpExpr(n.Target, tNext, true))
		sLine, sNext := node(next, false, "Start")
		if n.Start != nil {
			out.WriteString(sLine)
			out.WriteString(dumpExpr(n.Start, sNext, true))
		} else {
			out.WriteString(sLine + tNext + "<nil>\n")
		}
		eLine, eNext := node(next, true, "End")
		if n.End != nil {
			out.WriteString(eLine)
			out.WriteString(dumpExpr(n.End, eNext, true))
		} else {
			out.WriteString(eLine + eNext + "<nil>\n")
		}
		return out.String()
	default:
		line, _ := node(indent, last, fmt.Sprintf("<unknown %T>", n))
		return line
	}
}

func dumpType(t TypeNode, indent string, last bool) string {
	switch n := t.(type) {
	case *SimpleType:
		line, _ := node(indent, last, "Type "+n.Name)
		return line
	case *ListType:
		line, next := node(indent, last, "ListType")
		return line + dumpType(n.Element, next, true)
	case *FunctionType:
		line, next := node(indent, last, "FunctionType")
		var out strings.Builder
		out.WriteString(line)
		pLine, pNext := node(next, false, "Parameters")
		out.WriteString(pLine)
		for i, p := range n.Parameters {
			out.WriteString(dumpType(p, pNext, i == len(n.Parameters)-1))
		}
		rLine, rNext := node(next, true, "Return")
		out.WriteString(rLine)
		out.WriteString(dumpType(n.Return, rNext, true))
		return out.String()
	default:
		line, _ := node(indent, last, fmt.Sprintf("<unknown type %T>", n))
		return line
	}
}
