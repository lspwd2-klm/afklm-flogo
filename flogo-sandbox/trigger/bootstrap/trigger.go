package bootstrap

import (
	"context"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// MyTriggerFactory My Trigger factory
type MyTriggerFactory struct {
	metadata *trigger.Metadata
}

var log = logger.GetLogger("afklm-bootstrap-trigger")

// NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyTriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *MyTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &BootstrapTrigger{metadata: t.metadata, config: config}
}

// BootstrapTrigger is a stub for your Trigger implementation
type BootstrapTrigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	handlers []*trigger.Handler
}

// Initialize implements trigger.Init.Initialize
func (t *BootstrapTrigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	if t.handlers != nil {
		log.Info("Bootstrap trigger will start %d initialization flows", len(t.handlers))
	} else {
		log.Info("There are NO handlers attached to this trigger.")
	}
	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *BootstrapTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *BootstrapTrigger) Start() error {
	// start the trigger
	if t.handlers != nil {
		for _, h := range t.handlers {
			_, err := h.Handle(context.Background(), t.config.Settings)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to execute handler %s: %s", h.String(), err))
				return err
			}
		}
	} else {
		log.Warn("Bootstrap trigger wasn't able to find any handlers to trigger.")
	}

	return nil
}

// Stop implements trigger.Trigger.Start
func (t *BootstrapTrigger) Stop() error {
	// stop the trigger
	return nil
}
