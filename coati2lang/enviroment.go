package coati2lang

type Enviroment struct {
	Enclosing *Enviroment
	Values    map[string]interface{}
}

func NewEnviroment(enclosing *Enviroment) *Enviroment {
	return &Enviroment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}

func (e *Enviroment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Enviroment) Get(name string) (interface{}, bool) {
	if value, ok := e.Values[name]; ok {
		return value, ok
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, false
}

func (e *Enviroment) Assign(name string, value interface{}) {
	if _, ok := e.Values[name]; ok {
		e.Values[name] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

}
