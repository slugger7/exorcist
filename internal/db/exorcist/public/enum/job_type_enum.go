//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var JobTypeEnum = &struct {
	UpdateExistingVideos postgres.StringExpression
	ScanPath             postgres.StringExpression
	GenerateChecksum     postgres.StringExpression
}{
	UpdateExistingVideos: postgres.NewEnumValue("update_existing_videos"),
	ScanPath:             postgres.NewEnumValue("scan_path"),
	GenerateChecksum:     postgres.NewEnumValue("generate_checksum"),
}
