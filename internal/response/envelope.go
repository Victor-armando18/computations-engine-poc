package response

import (
	"github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"
	"github.com/dolphin-sistemas/computations-engine/core"
)

type Envelope struct {
	StateFragment canonical.OrderState `json:"stateFragment"`
	ServerDelta   interface{}          `json:"serverDelta,omitempty"`
	RulesVersion  string               `json:"rulesVersion"`
	Reasons       []core.Reason        `json:"reasons"`
	CorrelationID string               `json:"correlationId,omitempty"`
}

func Build(state canonical.OrderState, reasons []core.Reason) Envelope {
	return Envelope{
		StateFragment: state,
		ServerDelta:   nil,
		RulesVersion:  "2026.01",
		Reasons:       reasons,
		CorrelationID: "",
	}
}
