package runtime

import (
	"fmt"

	"github.com/mattn/anko/vm"
)

func (r *Runtime) Run(script string) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			switch v := recovered.(type) {
			case *AGSError:
				err = v
			case *AGSMultiError:
				err = v
			case error:
				err = fmt.Errorf("runtime panic: %w", v)
			default:
				err = fmt.Errorf("runtime panic: %v", recovered)
			}
		}
	}()

	_, err = vm.Execute(r.env, nil, script)
	if err != nil {
		return err
	}

	if err := r.validateSpec(); err != nil {
		return err
	}

	return nil
}
