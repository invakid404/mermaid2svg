package webdriver

import (
	"github.com/alexliesenfeld/health"
)

func Checker(driver *Driver) health.CheckerOption {
	return health.WithCheck(
		health.Check{
			Name:  "webdriver",
			Check: driver.Ping,
		},
	)
}
