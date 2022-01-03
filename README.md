# wmi
[![Go Reference](https://pkg.go.dev/badge/github.com/moonchant12/wmi)](https://pkg.go.dev/github.com/moonchant12/wmi)

Package wmi provides an interface to WMI. (Windows Management Instrumentation)

## Install
```cmd
go get -v github.com/moonchant12/wmi
```

## Import
```Go
import "github.com/moonchant12/wmi"
```

## Usage

### Call method with params
```Go
// Create process using WMI
func main() {
	// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-process
	// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/create-method-in-class-win32-process
	// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processstartup
	wmi.With("Win32_ProcessStartup", func(objStartupConfig *ole.IDispatch) error {
		var pid int32
		objStartupConfig.PutProperty("CreateFlags", windows.CREATE_NEW_CONSOLE)
		objStartupConfig.PutProperty("Title", "nice")
		log.Println(wmi.CallMethod("Win32_Process", "Create", "cmd /C notepad.exe", "C:\\Windows\\System32", objStartupConfig, &pid))
		log.Println(pid)
		return nil
	})
}
```
Outcome
```
2022/01/03 00:00:00 0 <nil>
2022/01/03 00:00:00 9996
```

### Query
```Go
// Query system processes information using WMI
func main() {
	var res []struct {
		Name        string
		HandleCount uint32
	}
	wmi.Query("SELECT Name, HandleCount FROM Win32_Process", &res)
	log.Println(res)
}
```
Outcome
```
2022/01/03 00:00:00 [{System Idle Process 0} {System 3860} {Secure System 0} {Registry 0}
...
```
