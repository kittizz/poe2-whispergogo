package main

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	kernel32                     = windows.NewLazySystemDLL("kernel32.dll")
	psapi                        = windows.NewLazySystemDLL("psapi.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procGetModuleFileNameExW     = psapi.NewProc("GetModuleFileNameExW")
)

func getForegroundWindowProcessName() string {
	hwnd, _, _ := procGetForegroundWindow.Call()

	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

	handle, _, _ := procOpenProcess.Call(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, 0, uintptr(pid))
	if handle == 0 {
		return ""
	}
	defer windows.CloseHandle(windows.Handle(handle))

	var buffer [260]uint16
	procGetModuleFileNameExW.Call(handle, 0, uintptr(unsafe.Pointer(&buffer[0])), 260)

	return syscall.UTF16ToString(buffer[:])
}

func main() {
	targetProcess := "PathOfExileSteam.exe"

	for {
		processPath := getForegroundWindowProcessName()
		processName := strings.ToLower(processPath[strings.LastIndex(processPath, "\\")+1:])

		if processName == strings.ToLower(targetProcess) {
			fmt.Printf("%s is currently focused\n", targetProcess)
		} else {
			fmt.Printf("%s is not focused (Current: %s)\n", targetProcess, processName)
		}

		time.Sleep(1 * time.Second)
	}
}
