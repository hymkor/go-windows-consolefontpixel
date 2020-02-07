package consolefontpixel

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var gdi32 = windows.NewLazySystemDLL("Gdi32.dll")
var procGetTextExtentPoint32 = gdi32.NewProc("GetTextExtentPoint32W")

var user32 = windows.NewLazySystemDLL("user32.dll")
var procGetDc = user32.NewProc("GetDC")
var procReleaseDc = user32.NewProc("GetDC")

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
var procGetCurrentConsoleFontEx = kernel32.NewProc("GetCurrentConsoleFontEx")

func getTextExtentPoint32(hdc uintptr, s string) (int, int, error) {
	utf16, err := windows.UTF16FromString(s)
	if err != nil {
		return -1, -1, err
	}
	c := len(utf16)
	p := &utf16[0]

	var size struct {
		cx int32 //long
		cy int32 //long
	}

	rc, _, err := procGetTextExtentPoint32.Call(hdc,
		uintptr(unsafe.Pointer(p)),
		uintptr(c),
		uintptr(unsafe.Pointer(&size)))
	if rc == 0 {
		return -1, -1, err
	}
	return int(size.cx), int(size.cy), nil
}

func getDc(hWnd uintptr) (uintptr, error) {
	hDc, _, _ := procGetDc.Call(hWnd)
	if hDc == 0 {
		return 0, errors.New("GetDC failed")
	}
	return hDc, nil
}

func releaseDc(hWnd, hDc uintptr) bool {
	isReleased, _, _ := procReleaseDc.Call(hWnd, hDc)
	return isReleased != 0
}

func getConsoleWindow() (uintptr, error) {
	hWnd, _, _ := procGetConsoleWindow.Call()
	if hWnd == 0 {
		return 0, errors.New("no such assciated console")
	}
	return hWnd, nil
}

func getCurrentConsoleFontEx(hConsoleOutput uintptr, maxWindows bool) (int, int, error) {
	var buffer struct {
		cbSize     uint32              // ULONG
		nFont      uint32              // DWORD
		dwFontSize windows.Coord       // COORD
		FontFamily uint                // UINT
		FontWeight uint                // UINT
		FaceName   [LF_FACESIZE]uint16 // WCHAR
	}
	var b uintptr = 0
	if maxWindows {
		b = 1
	}
	buffer.cbSize = uint32(unsafe.Sizeof(buffer))

	rc, _, err := procGetCurrentConsoleFontEx.Call(hConsoleOutput, b,
		uintptr(unsafe.Pointer(&buffer)))
	if rc == 0 {
		return 0, 0, err
	}
	return int(buffer.dwFontSize.X), int(buffer.dwFontSize.Y), nil
}

type DC struct {
	hWnd uintptr
	hDc  uintptr
}

func OpenDC() (*DC, error) {
	hWnd, err := getConsoleWindow()
	if err != nil {
		return nil, err
	}
	hDc, err := getDc(hWnd)
	if err != nil {
		return nil, err
	}
	return &DC{hWnd: hWnd, hDc: hDc}, nil
}

func (dc *DC) Close() {
	releaseDc(dc.hWnd, dc.hDc)
}

func (dc *DC) GetTextExtentPoint(s string) (int, int, error) {
	w, h, err := getTextExtentPoint32(dc.hDc, s)
	if err != nil {
		return -1, -1, fmt.Errorf("GetTextExtentPoint: %w", err)
	}
	return w, h, nil
}

func GetTextExtentPoint(s string) (int, int, error) {
	dc, err := OpenDC()
	if err != nil {
		return -1, -1, err
	}
	defer dc.Close()

	return dc.GetTextExtentPoint(s)
}

func GetCurrentConsoleFont() (int, int, error) {
	w, h, err := getCurrentConsoleFontEx(uintptr(windows.Stdout), false)
	if err != nil {
		return w, h, fmt.Errorf("GetCurrentConsoleFont: %w", err)
	}
	return w, h, nil
}
