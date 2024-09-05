package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ModelTestSuite is used for checking that the models in this package match
// the SQL `tables.sql` migration file. All model fields must be tagged with
// `table:"table_name"` where `table_name` is the name of the table in the
// database.
type ModelTestSuite struct {
	suite.Suite
}

func TestModel(t *testing.T) {

}
