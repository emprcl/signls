package main

import (
	"context"
	"cykl/midi"
	"cykl/sequencer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	seq := sequencer.New(midi)
	seq.TogglePlay()

	<-ctx.Done()
}
