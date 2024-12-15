package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// findProcessByName searches for a process by name and returns its path.
func findProcessByName(targetName string) (string, error) {
	// Snapshot all processes
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", fmt.Errorf("could not create process snapshot: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	// Iterate through the process list
	err = windows.Process32First(snapshot, &procEntry)
	if err != nil {
		return "", fmt.Errorf("could not get first process: %w", err)
	}

	for {
		processName := windows.UTF16ToString(procEntry.ExeFile[:])
		if strings.EqualFold(processName, targetName) {
			// Get executable path
			path, err := getProcessPath(procEntry.ProcessID)
			if err != nil {
				return "", fmt.Errorf("error getting process path: %w", err)
			}
			return path, nil // Found the desired process
		}

		// Move to the next process
		if err := windows.Process32Next(snapshot, &procEntry); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break // No more processes
			}
			return "", fmt.Errorf("error iterating processes: %w", err)
		}
	}
	return "", fmt.Errorf("process %s not found", targetName)
}

// getProcessPath retrieves the path of a process by its PID.
func getProcessPath(pid uint32) (string, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", fmt.Errorf("could not open process (PID: %d): %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	var buffer [windows.MAX_PATH]uint16
	size := uint32(len(buffer)) // Size of buffer in characters
	err = windows.QueryFullProcessImageName(handle, 0, &buffer[0], &size)
	if err != nil {
		return "", fmt.Errorf("could not get process image name: %w", err)
	}
	return windows.UTF16ToString(buffer[:]), nil
}

// tailFile continuously reads new lines from a file as they are written.
func tailFile(filePath string, quit chan os.Signal) error {
	// Open file in read-only mode
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Seek to the end of the file to ignore existing content
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("could not seek to end of file: %w", err)
	}

	reader := bufio.NewReader(file)
	done := make(chan struct{})

	// Goroutine to keep reading new lines from the file
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// Attempt to read a new line
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						// Wait a bit before re-checking for new lines
						time.Sleep(1 * time.Second)
					} else {
						fmt.Printf("Error reading file: %v\n", err)
						return
					}
				} else {
					// Print the new line read from the file
					fmt.Print(line)
				}
			}
		}
	}()

	// Wait for quit signal and handle shutdown
	<-quit
	fmt.Println("Shutting down gracefully...")
	close(done)
	return nil
}
