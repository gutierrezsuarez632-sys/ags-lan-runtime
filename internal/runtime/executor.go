package runtime

import "github.com/mattn/anko/vm"

func (r *Runtime) Run(script string) error {
    _, err := vm.Execute(r.env, nil, script)
    return err
}