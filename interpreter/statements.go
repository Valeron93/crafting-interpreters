package interpreter

import (
	"fmt"

	"github.com/Valeron93/crafting-interpreters/ast"
)

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	var value any
	var err error
	if stmt.Init != nil {
		value, err = i.Eval(stmt.Init)
		if err != nil {
			return nil, err
		}
	}
	i.env.Define(stmt.Name, value)
	return nil, nil
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) (any, error) {
	_, err := i.Eval(stmt.Expr)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	value, err := i.Eval(stmt.Expr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", value)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (any, error) {

	cond, err := i.Eval(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTrue(cond) {
		i.execute(stmt.Then)
	} else if stmt.Else != nil {
		i.execute(stmt.Else)
	}
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	subEnv := NewSubEnvironment(i.env)
	err := i.executeBlock(stmt.Statements, subEnv)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {

	for true {
		cond, err := i.Eval(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !i.isTrue(cond) {
			break
		}
		err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil

}
