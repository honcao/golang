package decorator

import (
	"net/http"
)

// Preparer to modify the webrequest
type Preparer interface {
	Prepare(*http.Request) (*http.Request, error)
}

// PreparerFunc implements Preparer
type PreparerFunc func(*http.Request) (*http.Request, error)

// Prepare implements Preparer interface
func (p PreparerFunc) Prepare(r *http.Request) (*http.Request, error) {
	return p(r)
}

// PreparerDecorator  contructor to help build the chain of decorators
type PreparerDecorator func(Preparer) Preparer

// CreatePreparer Create Preparer and decorators
func CreatePreparer(decorators ...PreparerDecorator) Preparer {
	return DecoratePreparer(Preparer(PreparerFunc(func(r *http.Request) (*http.Request, error) {
		return r, nil
	})), decorators...)
}

// DecoratePreparer decorate the Preparer
func DecoratePreparer(p Preparer, decorators ...PreparerDecorator) Preparer {
	for _, decorateor := range decorators {
		p = decorateor(p)
	}
	return p
}

// SetMethod udpate the method of http.Request
func SetMethod(method string) PreparerDecorator {
	return PreparerDecorator(func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Method = method
			return p.Prepare(r)
		})
	})
}
