package telebot

import "log"

// Pump is an interface representing a fuel pump.
type Pump interface {
	startFuelling()
	stopFuelling()
}

// FuelPump is a struct representing a fuel pump with a unique ID.
type FuelPump struct {
	pumpID uint64
}

// startFuelling starts the fueling process for the fuel pump.
func (p *FuelPump) startFuelling() {
	log.Printf("FuelPump #%d > started fuelling.", p.pumpID)
}

// stopFuelling stops the fueling process for the fuel pump.
func (p *FuelPump) stopFuelling() {
	log.Printf("FuelPump #%d > stopped fuelling.", p.pumpID)
}

// Command is an interface representing a generic command to be executed.
type Command interface {
	execute()
}

// StartCommand is a struct representing a command to start fueling for a pump.
type StartCommand struct {
	pump Pump
}

// execute starts the fueling process for the associated pump.
func (c *StartCommand) execute() {
	c.pump.startFuelling()
}

// StopCommand is a struct representing a command to stop fueling for a pump.
type StopCommand struct {
	pump Pump
}

// execute stops the fueling process for the associated pump.
func (c *StopCommand) execute() {
	c.pump.stopFuelling()
}
