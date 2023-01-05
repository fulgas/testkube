package testkube

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
)

type Tests []Test

func (t Tests) Table() (header []string, output [][]string) {
	header = []string{"Name", "Type", "Created", "Labels", "Schedule"}
	for _, e := range t {
		output = append(output, []string{
			e.Name,
			e.Type_,
			e.Created.String(),
			MapToString(e.Labels),
			e.Schedule,
		})
	}

	return
}

func (t Test) GetObjectRef(namespace string) *ObjectRef {
	return &ObjectRef{
		Name:      t.Name,
		Namespace: namespace,
	}
}

func PrepareExecutorArgs(binaryArgs []string) ([]string, error) {
	executorArgs := make([]string, 0)
	for _, arg := range binaryArgs {
		r := csv.NewReader(strings.NewReader(arg))
		r.Comma = ' '
		r.LazyQuotes = true
		r.TrimLeadingSpace = true

		records, err := r.ReadAll()
		if err != nil {
			return nil, err
		}

		if len(records) != 1 {
			return nil, errors.New("single string expected")
		}

		executorArgs = append(executorArgs, records[0]...)
	}

	return executorArgs, nil
}

func (test *Test) QuoteTestTextFields() {
	if test.Content != nil && test.Content.Data != "" {
		test.Content.Data = fmt.Sprintf("%q", test.Content.Data)
	}

	if test.ExecutionRequest != nil {
		var fields = []*string{
			&test.ExecutionRequest.VariablesFile,
			&test.ExecutionRequest.JobTemplate,
			&test.ExecutionRequest.PreRunScript,
			&test.ExecutionRequest.ScraperTemplate,
		}

		for _, field := range fields {
			if *field != "" {
				*field = fmt.Sprintf("%q", *field)
			}
		}
	}
}
