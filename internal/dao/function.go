package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao的方法，dao作为接收者
func (d RuleDao) CreateFunction(Type string, threshold float64, description string) error {
	function := model.Function{Type: Type, Threshold: threshold, Description: description}
	return function.Create(d.engine)
}

func (d RuleDao) DeleteFunction(id uint32) error {
	function := model.Function{Id: id}
	return function.Delete(d.engine)
}

func (d RuleDao) GetFunction(id uint32) (model.Function, error) {
	function := model.Function{Id: id}
	return function.Get(d.engine)
}

func (d RuleDao) UpdateFunction(id uint32, Type string, threshold float64, description string) error {
	function := model.Function{Id: id, Type: Type, Threshold: threshold, Description: description}
	return function.Update(d.engine)
}

func (d RuleDao) ListFunctions() ([]model.Function, error) {
	function := model.Function{}
	return function.List(d.engine)
}
