package interpreter

import (
	"errors"

	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/scanner"
	"github.com/Valeron93/crafting-interpreters/util"
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
		return nil, util.ReportErrorOnToken(expr.Operator, "operands '%v' and '%v' are not compatible with binary operator '%v'", left, right, expr.Operator.Lexeme)
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
		return nil, util.ReportErrorOnToken(expr.Paren, "'%#v' is not callable", callee)
	}
	arity, varArg := f.Arity()
	if !varArg && len(args) != arity {
		return nil, util.ReportErrorOnToken(expr.Paren, "function '%v' expects %v arguments, got: %v", callee, arity, len(args))
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

func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) (any, error) {
	object, err := i.Eval(expr.Object)
	if err != nil {
		return nil, err
	}

	if object, ok := object.(Object); ok {
		return object.Get(expr.Name)
	}

	return nil, util.ReportErrorOnToken(expr.Name, "only classes and class instances have properties")
}

func (i *Interpreter) VisitSetExpr(expr *ast.SetExpr) (any, error) {
	object, err := i.Eval(expr.Object)
	if err != nil {
		return nil, err
	}

	if obj, ok := object.(Object); ok {
		value, err := i.Eval(expr.Value)
		if err != nil {
			return nil, err
		}

		err = obj.Set(expr.Name, value)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	return nil, util.ReportErrorOnToken(expr.Name, "only class instances have properties ")
}

func (i *Interpreter) VisitThisExpr(expr *ast.ThisExpr) (any, error) {
	return i.lookUpVar(expr.Keyword, expr)
}

func (i *Interpreter) VisitGetKeyExpr(expr *ast.GetKeyExpr) (any, error) {
	object, err := i.Eval(expr.Object)
	if err != nil {
		return nil, err
	}

	if obj, ok := object.(IndexableObject); ok {

		key, err := i.Eval(expr.Key)
		if err != nil {
			return nil, err
		}

		value, err := obj.GetKeyValue(i, expr.Bracket, key)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	return nil, util.ReportErrorOnToken(expr.Bracket, "indexing is only supported on objects")

}

func (i *Interpreter) VisitSetKeyExpr(expr *ast.SetKeyExpr) (any, error) {

	object, err := i.Eval(expr.Object)
	if err != nil {
		return nil, err
	}

	if obj, ok := object.(IndexableObject); ok {
		key, err := i.Eval(expr.Key)
		if err != nil {
			return nil, err
		}

		value, err := i.Eval(expr.Value)
		if err != nil {
			return nil, err
		}

		return nil, obj.SetKeyValue(i, expr.Bracket, key, value)
	}
	return nil, util.ReportErrorOnToken(expr.Bracket, "indexing is not supported on this object")

}

func (i *Interpreter) VisitSuperExpr(expr *ast.SuperExpr) (any, error) {
	distance := i.locals[expr]

	super, err := i.env.GetAt(distance, expr.Keyword)
	if err != nil {
		return nil, err
	}

	superclass, ok := super.(*Class)
	if !ok {
		return nil, util.ReportErrorOnToken(expr.Keyword, "failed to find superclass")
	}

	// copy the keyword so we can get location info
	thisKeyword := expr.Keyword
	thisKeyword.Lexeme = "this"
	object, err := i.env.GetAt(distance-1, thisKeyword)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*ClassInstance)
	if !ok {
		return nil, util.ReportErrorOnToken(expr.Keyword, "failed to obtain instance")
	}

	method, ok := superclass.FindMethod(expr.Method.Lexeme)
	if !ok {
		return nil, util.ReportErrorOnToken(expr.Keyword, "'super' has no '%v' method", expr.Method.Lexeme)
	}

	return method.Bind(instance), nil

}
