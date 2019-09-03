package fuel

import (
	"errors"
	"time"

	"github.com/rightjoin/dorm"
	"github.com/rightjoin/fig"
)

type HealthService struct {
	Service      `root:"-"`
	HealthChecks func() []HealthStatus

	check GET `route:"health-check" middleware:"-" cache:"-" ttl:"-" wrap:"true"`
}

func (h *HealthService) Check() ([]HealthStatus, error) {

	out := make([]HealthStatus, 0)

	checks := fig.StringSliceOr(nil, "auto-heath-check")
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
	var err error = nil
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

type HealthStatus struct {
	Name     string    `json:"name"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TestedAt time.Time `json:"tested_at"`
}
