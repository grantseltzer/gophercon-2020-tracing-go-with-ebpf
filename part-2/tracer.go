package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/iovisor/gobpf/bcc"
)

const eBPFProgram = `
#include <uapi/linux/ptrace.h>
#include <linux/sched.h>

BPF_PERF_OUTPUT(events);

int function_was_called(struct pt_regs *ctx) {
	char message[] = "handler function was called!";
	events.perf_submit(ctx, &message, sizeof(message));
	return 0;
}
`

func main() {
	bpfModule := bcc.NewModule(eBPFProgram, []string{})

	uprobeFD, err := bpfModule.LoadUprobe("function_was_called")
	if err != nil {
		log.Fatal(err)
	}

	err = bpfModule.AttachUprobe(os.Args[1], "main.handlerFunction", uprobeFD, -1)
	if err != nil {
		log.Fatal(err)
	}

	table := bcc.NewTable(bpfModule.TableId("events"), bpfModule)
	channel := make(chan []byte)
	perfMap, err := bcc.InitPerfMap(table, channel, nil)
	if err != nil {
		log.Fatal(err)
	}

	perfMap.Start()
	defer perfMap.Stop()

	go func() {
		for {
			value := <-channel
			fmt.Println(string(value))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
