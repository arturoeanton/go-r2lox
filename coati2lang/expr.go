package coati2lang

type Expr interface {
	AcceptExpr(visitor Visitor) interface{}
}

type Stmt interface {
	AcceptStmt(visitor Visitor) interface{}
}

type Visitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitGroupingABSExpr(expr GroupingABS) interface{}
	VisitLiteralExpr(expr Literal) interface{}
	VisitUnaryExpr(expr Unary) interface{}
	VisitAssignExpr(stmt Assign) interface{}
	VisitLogicalExpr(expr Logical) interface{}
	VisitCallExpr(expr Call) interface{}

	VisitExpressionStmt(stmt Expression) interface{}
	VisitVar(stmt Var) interface{}

	VisitVariableExpr(expr Var) interface{}
	VisitBlockStmt(stmt Block) interface{}
	VisitIfStmt(stmt If) interface{}
	VisitWhileStmt(stmt While) interface{}
	VisitFunctionStmt(stmt Function) interface{}
	VisitReturnStmt(stmt Return) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

type GroupingABS struct {
	Expression Expr
}

func (g GroupingABS) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitGroupingABSExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Operator Token
	Value    Expr
}

func (u Unary) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

type Expression struct {
	Expression Expr
}

func (e Expression) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

type Condition struct {
	Condition Expr
}

func (c Condition) AcceptStmt(visitor Visitor) interface{} {
	return visitor
}

/*type Print struct {
	Expression Expr
}

func (p Print) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitPrint(p)
}*/

type Var struct {
	Name             Token
	InitializerVal   Expr
	InitializerArray []Expr
	InitializerMap   []ItemVar
	InitializerFx    Stmt
	Selectors        [][]Expr
	Sub              bool
	SizeArrayInit    int
}
type ItemVar struct {
	Key   Expr
	Value Expr
}

func (v Var) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

func (v Var) AcceptStmt(visitor Visitor) interface{} {

	return visitor.VisitVar(v)
}

type Assign struct {
	Name      Token
	Value     Expr
	Selectors [][]Expr
	Sub       bool
}

func (a Assign) AcceptExpr(visitor Visitor) interface{} {

	return visitor.VisitAssignExpr(a)
}

type Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
	This      Var
	SubCall   *Call
}

func (c Call) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitCallExpr(c)
}

type Super struct {
	Keyword Token

	Method Token
}

func (s Super) AcceptExpr(visitor Visitor) interface{} {
	return nil
}

type Class struct {
	Name       Token
	Methods    []Function
	Superclass Var
}

func (c Class) AcceptExpr(visitor Visitor) interface{} {
	return nil
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i If) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitIfStmt(i)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (w While) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitWhileStmt(w)
}

type Block struct {
	Statements []Stmt
}

func (b Block) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

type Return struct {
	Keyword Token
	Value   Expr
	Result  interface{}
}

func (r Return) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitReturnStmt(r)
}

type Break struct {
	Keyword Token
}

func (b Break) AcceptExpr(visitor Visitor) interface{} {
	return nil
}

type Continue struct {
	Keyword Token
}

func (c Continue) AcceptExpr(visitor Visitor) interface{} {
	return nil
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l Logical) AcceptExpr(visitor Visitor) interface{} {
	return visitor.VisitLogicalExpr(l)

}

type Function struct {
	Name       Token
	Parameters []Token
	Body       []Stmt
	Closure    *Enviroment
}

func (f Function) AcceptStmt(visitor Visitor) interface{} {
	return visitor.VisitFunctionStmt(f)
}

func (f Function) Call(i *Interpreter, arguments []interface{}, this interface{}) interface{} {

	enviroment := NewEnviroment(f.Closure)
	for i, param := range f.Parameters {
		enviroment.Define(param.Lexeme, arguments[i])
	}
	enviroment.Define("this", this)

	var value interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				return_result, ok := r.(Return)
				if ok {
					value = return_result.Result
					return
				}

				if r == "break" {
					return
				}
				if r == "continue" {
					return
				}
				panic(r)
			}
		}()
		value = i.executeBlock(f.Body, *enviroment)
	}()

	return value
}

func (f Function) Arity() int {
	return len(f.Parameters)
}
