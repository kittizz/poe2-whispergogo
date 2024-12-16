package main

import (
	"fmt"
	"math/rand"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")

	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procEnumWindows              = user32.NewProc("EnumWindows")

	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func findProcess(name string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	err = windows.Process32First(snapshot, &pe)
	if err != nil {
		return 0, err
	}

	for {
		if strings.EqualFold(windows.UTF16ToString(pe.ExeFile[:]), name) {
			return pe.ProcessID, nil
		}
		err = windows.Process32Next(snapshot, &pe)
		if err != nil {
			break
		}
	}
	return 0, fmt.Errorf("process not found")
}

func setGameWindowFocus(pid uint32) bool {
	var targetHWND windows.Handle

	// callback function สำหรับ EnumWindows
	callback := func(hwnd windows.Handle, lparam uintptr) uintptr {
		var currentPID uint32
		procGetWindowThreadProcessId.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&currentPID)),
		)

		if currentPID == pid {
			targetHWND = hwnd
			procSetForegroundWindow.Call(uintptr(hwnd))
			return 0 // หยุดการค้นหา
		}
		return 1 // ค้นหาต่อ
	}

	// เรียกใช้ EnumWindows ผ่าน syscall
	procEnumWindows.Call(
		syscall.NewCallback(callback),
		0,
	)

	return targetHWND != 0
}

func main() {
	processName := "PathOfExileSteam.exe"

	for {
		pid, err := findProcess(processName)
		if err != nil {
			fmt.Printf("%s not found. Waiting...\n", processName)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Found %s (PID: %d). Applying fake gameplay...\n", processName, pid)

		for {
			if setGameWindowFocus(pid) {
				fmt.Println("Successfully focused game window")
			} else {
				fmt.Println("Failed to focus game window")
			}

			time.Sleep(time.Second * time.Duration(rnd.Intn(3)+1))

			_, err := findProcess(processName)
			if err != nil {
				fmt.Printf("%s is no longer running. Searching again...\n", processName)
				break
			}
		}
	}
}
