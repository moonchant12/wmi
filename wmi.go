//go:build windows

/*
Package wmi provides an interface to WMI. (Windows Management Instrumentation)
*/
package wmi

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/yusufpapurcu/wmi"
)

// With usage: wmi.With("Win32_ProcessStartup", func(objStartupConfig *ole.IDispatch) { ... })
func With(className string, do func(class *ole.IDispatch) error) error {
	// Initialize COM
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	defer ole.CoUninitialize()

	// Connect to WMI through the IWbemLocator::ConnectServer method
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return fmt.Errorf("With has failed to initialize 1: %v", err)
	}
	defer unknown.Release()

	// convert it to IDispatch
	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("With has failed to initialize 2: %v", err)
	}
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return fmt.Errorf("With has failed to initialize 3: %v", err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// CoInit success

	// Get class
	classRaw, err := oleutil.CallMethod(service, "Get", className)
	if err != nil {
		return fmt.Errorf("With class %s: %v", className, err)
	}
	class := classRaw.ToIDispatch()
	defer classRaw.Clear()

	// Do
	if do != nil {
		return do(class)
	}

	return nil
}

// CallMethod usage: wmi.CallMethod("Win32_Process", "Create", "cmd /C notepad.exe", "C:\\Windows\\System32", objStartupConfig, &pid)
//
// CallMethod calls a method named methodName on an instance of the class named
// className, with the given params.
//
// CallMethod is a wrapper around DefaultClient.CallMethod.
func CallMethod(className, methodName string, params ...interface{}) (int32, error) {
	return wmi.CallMethod(nil, className, methodName, params)
}

// Query usage: wmi.Query("SELECT Name, HandleCount FROM Win32_Process", &res)
//
// var res []struct {
// Name        string
// HandleCount uint32
// }
//
// Query runs the WQL query and appends the values to dst.
//
// dst must have type *[]S or *[]*S, for some struct type S. Fields selected in
// the query must have the same name in dst. Supported types are all signed and
// unsigned integers, time.Time, string, bool, or a pointer to one of those.
// Array types are not supported.
//
// By default, the local machine and default namespace are used. These can be
// changed using connectServerArgs. See
// https://docs.microsoft.com/en-us/windows/desktop/WmiSdk/swbemlocator-connectserver
// for details.
//
// Query is a wrapper around DefaultClient.Query.
func Query(query string, dst interface{}, connectServerArgs ...interface{}) error {
	return wmi.Query(query, dst, connectServerArgs...)
}
