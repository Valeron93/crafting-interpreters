package interpreter

import (
	"errors"
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
)

func (i *Interpreter) VisitVarExpr(expr *ast.VarExpr) (any, error) {
	return i.lookUpVar(expr.Name, expr)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	value, err := i.Eval(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[expr]
	if ok {
		err = i.env.AssignAt(distance, expr.Name, value)
		if err != nil {
			return nil, err
		}
	} else {
		err = i.globals.Assign(expr.Name, value)
		if err != nil {
			return nil, err
		}
	}

	err = i.env.Assign(expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	left, err := i.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.EqualEqual:
		return left == right, nil
	case scanner.BangEqual:
		return left != right, nil
	}
	lhsFloat, lhsIsFloat := left.(float64)
	rhsFloat, rhsIsFloat := right.(float64)

	lhsString, lhsIsString := left.(string)
	rhsString, rhsIsString := right.(string)

	if lhsIsFloat && rhsIsFloat {

		switch expr.Operator.Type {
		case scanner.Plus, scanner.Minus, scanner.Star, scanner.Slash:
			return floatOperator(expr.Operator, lhsFloat, rhsFloat), nil

		case scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual:
			return floatLogicOperator(expr.Operator, lhsFloat, rhsFloat), nil
		}
	} else if (lhsIsString && rhsIsString) && expr.Operator.Type == scanner.Plus {
		return lhsString + rhsString, nil
	} else {
		return nil, fmt.Errorf("operands '%v' and '%v' are not compatible with binary operator '%v'", left, right, expr.Operator.Lexeme)
	}

	return nil, errors.ErrUnsupported
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.Eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.Bang:
		return !(i.isTrue(right)), nil
	case scanner.Minus:
		return -(right.(float64)), nil
	}
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.Eval(expr.Expression)
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	left, err := i.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Type == scanner.Or {
		if i.isTrue(left) {
			return left, nil
		}
	} else {
		if !i.isTrue(left) {
			return left, nil
		}
	}
	return i.Eval(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (any, error) {

	callee, err := i.Eval(expr.Callee)
	if err != nil {
		return nil, err
	}

	args := make([]any, 0)

	for _, arg := range expr.Args {
		value, err := i.Eval(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}

	f, ok := callee.(Callable)
	if !ok {
		return nil, fmt.Errorf("`%#v` is not callable", callee)
	}
	arity, varArg := f.Arity()
	if !varArg && len(args) != arity {
		return nil, fmt.Errorf("function expects %v arguments, got: %v", arity, len(args))
	}

	return f.Call(i, args)
}

func (i *Interpreter) VisitLambdaExpr(expr *ast.LambdaExpr) (any, error) {
	return &CallableObject{
		Declaration: &ast.FuncDeclStmt{
			Name: scanner.Token{
				Type:   scanner.Ident,
				Lexeme: "lambda",
			},
			Params: expr.Params,
			Body:   expr.Body,
		},
		Closure: i.env,
	}, nil
}
