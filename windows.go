package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	processCheckInterval = 5 * time.Second
	fileReadInterval     = 1 * time.Second
)

type ProcessWatcher interface {
	WatchProcess(ctx context.Context, gameStatus *bool) error
	GetProcessPath(pid uint32) (string, error)
}

type FileWatcher interface {
	TailFile(ctx context.Context, filePath string) error
}

type DefaultProcessWatcher struct {
	wg         sync.WaitGroup
	lineStream chan string
}

type DefaultFileWatcher struct {
	lineStream chan string
}

func NewProcessWatcher(lineStream chan string) ProcessWatcher {
	return &DefaultProcessWatcher{
		lineStream: lineStream,
	}
}

func NewFileWatcher(lineStream chan string) FileWatcher {
	return &DefaultFileWatcher{
		lineStream: lineStream,
	}
}

func (pw *DefaultProcessWatcher) WatchProcess(ctx context.Context, isProcessRunning *bool) error {
	fw := NewFileWatcher(pw.lineStream)

	var fileWatcherCancel context.CancelFunc

	pw.wg.Add(1)
	go func() {
		defer pw.wg.Done()
		for {
			select {
			case <-ctx.Done():
				if fileWatcherCancel != nil {
					fileWatcherCancel()
				}
				return
			default:
				clientPath, err := pw.findProcessByName()
				if err == nil && !*isProcessRunning {
					*isProcessRunning = true
					pw.wg.Add(1)

					var fileCtx context.Context
					fileCtx, fileWatcherCancel = context.WithCancel(ctx)

					fmt.Printf("Found PoE Client.txt Path: %s\n", clientPath)

					go func(path string) {
						defer pw.wg.Done()
						defer func() {
							*isProcessRunning = false
							if fileWatcherCancel != nil {
								fileWatcherCancel()
								fileWatcherCancel = nil
							}
						}()

						fmt.Println("Starting to tail Client.txt...")
						if err := fw.TailFile(fileCtx, path); err != nil {
							if err != context.Canceled {
								fmt.Printf("Error tailing file: %v\n", err)
							}
						}
					}(clientPath)

				} else if err != nil && *isProcessRunning {
					fmt.Println("PoE process stopped, waiting for restart...")
					if fileWatcherCancel != nil {
						fileWatcherCancel()
						fileWatcherCancel = nil
					}
					*isProcessRunning = false
				}
				time.Sleep(processCheckInterval)
			}
		}
	}()

	return nil
}

func (fw *DefaultFileWatcher) TailFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("seek to end: %w", err)
	}

	reader := bufio.NewReader(file)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping file tail operation...")
			return ctx.Err()
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						time.Sleep(fileReadInterval)
						continue
					}
				}
				return fmt.Errorf("read file: %w", err)
			}
			select {
			case fw.lineStream <- strings.TrimSpace(line):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (pw *DefaultProcessWatcher) findProcessByName() (string, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", fmt.Errorf("create process snapshot: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err := windows.Process32First(snapshot, &procEntry); err != nil {
		return "", fmt.Errorf("get first process: %w", err)
	}

	for {
		processName := windows.UTF16ToString(procEntry.ExeFile[:])
		if strings.EqualFold(processName, POE2_PROCESS_NAME) {
			path, err := pw.GetProcessPath(procEntry.ProcessID)
			if err != nil {
				return "", fmt.Errorf("get process path: %w", err)
			}

			clientFilePath := filepath.Join(filepath.Dir(path), "logs", "Client.txt")
			return clientFilePath, nil
		}

		if err := windows.Process32Next(snapshot, &procEntry); err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return "", fmt.Errorf("iterate processes: %w", err)
		}
	}
	return "", fmt.Errorf("PoE process not found")
}

func (pw *DefaultProcessWatcher) GetProcessPath(pid uint32) (string, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", fmt.Errorf("open process (PID: %d): %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	var buffer [windows.MAX_PATH]uint16
	size := uint32(len(buffer))
	if err := windows.QueryFullProcessImageName(handle, 0, &buffer[0], &size); err != nil {
		return "", fmt.Errorf("get process image name: %w", err)
	}
	return windows.UTF16ToString(buffer[:]), nil
}
