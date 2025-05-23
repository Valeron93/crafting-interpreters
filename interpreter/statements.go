package interpreter

import (
	"github.com/Valeron93/crafting-interpreters/ast"
	"github.com/Valeron93/crafting-interpreters/util"
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
	err = i.env.Define(stmt.Name.Lexeme, value)
	return nil, err
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) (any, error) {
	_, err := i.Eval(stmt.Expr)
	return nil, err
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (any, error) {

	cond, err := i.Eval(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTrue(cond) {
		if err = i.execute(stmt.Then); err != nil {
			return nil, err
		}
	} else if stmt.Else != nil {
		if err = i.execute(stmt.Else); err != nil {
			return nil, err
		}
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

func (i *Interpreter) VisitFuncDeclStmt(stmt *ast.FuncDeclStmt) (any, error) {
	f := &CallableObject{
		Declaration: stmt,
		Closure:     NewSubEnvironment(i.env),
	}

	err := i.env.Define(stmt.Name.Lexeme, f)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) (any, error) {
	var value any
	var err error
	if stmt.Value != nil {
		value, err = i.Eval(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, &FunctionReturn{
		Value: value,
	}
}

func (i *Interpreter) VisitClassDeclStmt(stmt *ast.ClassDeclStmt) (any, error) {
	var superclass *Class
	if stmt.Superclass != nil {
		class, err := i.Eval(stmt.Superclass)
		if err != nil {
			return nil, err
		}

		if class, ok := class.(*Class); ok {
			superclass = class
		} else {
			return nil, util.ReportErrorOnToken(stmt.Superclass.Name, "'%v' is not a class", stmt.Superclass.Name.Lexeme)
		}
	}
	i.env.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.env = NewSubEnvironment(i.env)
		i.env.Define("super", superclass)
	}

	methods := make(map[string]ClassMethod)
	var init Callable
	for _, method := range stmt.Methods {
		if method.Func.Name.Lexeme == "init" {
			init = &CallableObject{
				Declaration: method.Func,
				Closure:     i.env,
			}
			continue
		}
		methods[method.Func.Name.Lexeme] = ClassMethod{
			Callable: &CallableObject{
				Declaration: method.Func,
				Closure:     i.env,
			},
			Static: method.Static,
		}
	}
	class := &Class{
		Name:        stmt.Name.Lexeme,
		Methods:     methods,
		Constructor: init,
		Superclass:  superclass,
	}

	if superclass != nil {
		i.env = i.env.enclosing
	}

	i.env.Assign(stmt.Name, class)
	return nil, nil
}

func (i *Interpreter) VisitMethodDeclStmt(stmt *ast.MethodDeclStmt) (any, error) {
	// we only need this function only to implement the StmtVisitor interface,
	// since MethodDeclStmt are interpreted in the VisitClassDeclStmt
	// and are not visited

	// TODO: maybe, there is a better approach?
	return nil, nil
}
