//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var VideoImageTypeEnum = &struct {
	Thumbnail postgres.StringExpression
	Chapter   postgres.StringExpression
}{
	Thumbnail: postgres.NewEnumValue("thumbnail"),
	Chapter:   postgres.NewEnumValue("chapter"),
}
