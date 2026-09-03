package tests

import (
	"github.com/goravel/framework/testing"

	"my_vehicles_an_backup_agent/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
