//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var JobStatus = &struct {
	NotStarted postgres.StringExpression
	InProgress postgres.StringExpression
	Failed     postgres.StringExpression
	Completed  postgres.StringExpression
	Cancelled  postgres.StringExpression
}{
	NotStarted: postgres.NewEnumValue("not_started"),
	InProgress: postgres.NewEnumValue("in_progress"),
	Failed:     postgres.NewEnumValue("failed"),
	Completed:  postgres.NewEnumValue("completed"),
	Cancelled:  postgres.NewEnumValue("cancelled"),
}
