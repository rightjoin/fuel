package fuel

import (
	"errors"
	"fmt"
	"time"

	"github.com/rightjoin/dorm"
	"github.com/rightjoin/fig"
)

type HealthService struct {
	Service      `root:"-"`
	HealthChecks func() []HealthStatus

	check GET `route:"health-check" middleware:"-" cache:"-" ttl:"-" wrap:"true"`
}

func (h *HealthService) Check() (out []HealthStatus, err error) {

	defer Recover(&out, &err)

	// Support both 'health-checks' and older 'auto-heath-check'
	var checks = []string{}
	{
		checks = fig.StringSliceOr(nil, "auto-heath-check")
		if checks == nil {
			checks = fig.StringSliceOr(nil, "health-checks")
		}
	}

	if checks != nil {
		for _, config := range checks {
			engine := fig.String(config, "engine")
			switch engine {
			case "mysql":
				dbo := dorm.GetORM(true)
				if errdb := dbo.DB().Ping(); errdb != nil {
					out = append(out, HealthStatus{
						Name:     config,
						Success:  false,
						Message:  errdb.Error(),
						TestedAt: time.Now(),
					})
				} else {
					out = append(out, HealthStatus{
						Name:     config,
						Success:  true,
						Message:  "ping successful",
						TestedAt: time.Now(),
					})
				}
			default:
				out = append(out, HealthStatus{
					Name:     config,
					Success:  false,
					Message:  "health check not implemented : " + engine,
					TestedAt: time.Now(),
				})
			}
		}
	}

	if h.HealthChecks != nil {
		out = append(out, h.HealthChecks()...)
	}

	// Is there any error?

	for _, h := range out {
		if h.Success == false {
			msg := h.Name
			if len(msg) > 0 {
				msg += " : "
			}
			msg += h.Message
			err = errors.New(msg)
		}
	}

	return out, err
}

// Recover will handle any panic.
func Recover(h *[]HealthStatus, errors *error) {

	var err error
	if r := recover(); r != nil {
		if err1, ok := r.(error); ok {

			err = err1
		}
		if err == nil {
			err = fmt.Errorf("%v", r)
		}

		*h = append(*h, HealthStatus{
			Name:     "Health-Check Service",
			Success:  false,
			Message:  "Cannot establish connection to db",
			TestedAt: time.Now(),
		})
		errors = &err
	}

}

type HealthStatus struct {
	Name     string    `json:"name"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TestedAt time.Time `json:"tested_at"`
}
