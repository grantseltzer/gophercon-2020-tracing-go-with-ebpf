package main

import (
	"encoding/binary"
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
	void* stackAddr = (void*)ctx->sp;

	long newArgument1 = 7;
	bpf_probe_write_user(stackAddr+8, &newArgument1, sizeof(newArgument1));
	
	return 0;
}
`

func main() {
	bpfModule := bcc.NewModule(eBPFProgram, []string{})

	uprobeFD, err := bpfModule.LoadUprobe("function_was_called")
	if err != nil {
		log.Fatal(err)
	}

	err = bpfModule.AttachUprobe(os.Args[1], "main.addTwoNumbers", uprobeFD, -1)
	if err != nil {
		log.Fatal(err)
	}

	table := bcc.NewTable(bpfModule.TableId("events"), bpfModule)
	perfChannel := make(chan []byte)
	perfMap, err := bcc.InitPerfMap(table, perfChannel, nil)
	if err != nil {
		log.Fatal(err)
	}

	perfMap.Start()
	defer perfMap.Stop()

	go func() {
		for {
			firstBytes := <-perfChannel
			first := binary.LittleEndian.Uint64(firstBytes)
			secondBytes := <-perfChannel
			second := binary.LittleEndian.Uint32(secondBytes)
			fmt.Printf("Arg1: %d\tArg2: %d\n", first, second)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
