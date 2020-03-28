package errors

import "encoding/json"

func NewMultiErrors() *MultiErrors {
	return &MultiErrors{
		errors: []error{},
	}
}

type MultiErrors struct {
	errors []error
}

func (this *MultiErrors) Add(errors ...error) *MultiErrors {
	for _, err := range errors {
		if err != nil {
			this.errors = append(this.errors, err)
		}
	}
	return this
}

func (this *MultiErrors) Error() string {
	msgs := make([]string, len(this.errors))
	for i := range msgs {
		msgs[i] = this.errors[i].Error()
	}
	bs, _ := json.Marshal(msgs)
	return string(bs)
}

func (this *MultiErrors) HasError() bool {
	return len(this.errors) > 0
}

func (this *MultiErrors) GetErrors() []error {
	return this.errors
}
