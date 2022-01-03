package wmi_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-ole/go-ole"
	"github.com/moonchant12/wmi"
	"golang.org/x/sys/windows"
)

func TestCreateProcess(t *testing.T) {
	// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-process
	// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/create-method-in-class-win32-process
	// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processstartup
	if errWith := wmi.With("Win32_ProcessStartup", func(objStartupConfig *ole.IDispatch) error {
		var pid int32
		if _, errPutProperty1 := objStartupConfig.PutProperty("CreateFlags", windows.CREATE_NEW_CONSOLE); errPutProperty1 != nil {
			t.Errorf("PutProperty 1: %v", errPutProperty1)
			return errPutProperty1
		}
		if _, errPutProperty2 := objStartupConfig.PutProperty("Title", "nice"); errPutProperty2 != nil {
			t.Errorf("PutProperty 2: %v", errPutProperty2)
			return errPutProperty2
		}
		ret, errCallMethod := wmi.CallMethod("Win32_Process", "Create", "cmd /C notepad.exe", "C:\\Windows\\System32", objStartupConfig, &pid)
		if errCallMethod != nil {
			t.Errorf("CallMethod: %v", errCallMethod)
			return errCallMethod
		}
		if ret != 0 {
			t.Errorf("Return code: %v", ret)
			return fmt.Errorf("Return code: %v", ret)
		}
		if pid == 0 {
			t.Errorf("Process ID: %v", pid)
			return fmt.Errorf("Process ID: %v", pid)
		}
		return nil
	}); errWith != nil {
		t.Errorf("With: %v", errWith)
	}

	// p, errProcess := os.StartProcess(
	// 	"notepad.exe",
	// 	[]string{},
	// 	&os.ProcAttr{
	// 		Dir: "C:\\Windows\\System32",
	// 		Env: os.Environ(),
	// 		Files: []*os.File{
	// 			os.Stdin,
	// 			os.Stdout,
	// 			os.Stderr,
	// 		},
	// 		Sys: nil,
	// 	},
	// )
	// if errProcess != nil {
	// 	log.Println("Failed to start process:", p, errProcess)
	// 	return
	// }
	// p.Release()

	// exec.Command("cmd.exe", "/C", "notepad.exe").Run()
	// exec.Command("cmd.exe", "/C", "notepad.exe").Start()
}

func TestQuery(t *testing.T) {
	var res []struct {
		Name        string
		HandleCount uint32
	}
	if errQuery := wmi.Query("SELECT Name, HandleCount FROM Win32_Process", &res); errQuery != nil {
		t.Errorf("Query: %v", errQuery)
	}
	log.Println(res)
}
