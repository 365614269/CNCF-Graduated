//go:build windows

// Code generated by 'go generate' using "github.com/Microsoft/go-winio/tools/mkwinsyscall"; DO NOT EDIT.

package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	return e
}

var (
	modadvapi32   = windows.NewLazySystemDLL("advapi32.dll")
	modbindfltapi = windows.NewLazySystemDLL("bindfltapi.dll")
	modcfgmgr32   = windows.NewLazySystemDLL("cfgmgr32.dll")
	modcimfs      = windows.NewLazySystemDLL("cimfs.dll")
	modiphlpapi   = windows.NewLazySystemDLL("iphlpapi.dll")
	modkernel32   = windows.NewLazySystemDLL("kernel32.dll")
	modnetapi32   = windows.NewLazySystemDLL("netapi32.dll")
	modntdll      = windows.NewLazySystemDLL("ntdll.dll")
	modoffreg     = windows.NewLazySystemDLL("offreg.dll")

	procLogonUserW                             = modadvapi32.NewProc("LogonUserW")
	procBfSetupFilter                          = modbindfltapi.NewProc("BfSetupFilter")
	procCM_Get_DevNode_PropertyW               = modcfgmgr32.NewProc("CM_Get_DevNode_PropertyW")
	procCM_Get_Device_ID_ListA                 = modcfgmgr32.NewProc("CM_Get_Device_ID_ListA")
	procCM_Get_Device_ID_List_SizeA            = modcfgmgr32.NewProc("CM_Get_Device_ID_List_SizeA")
	procCM_Get_Device_Interface_ListW          = modcfgmgr32.NewProc("CM_Get_Device_Interface_ListW")
	procCM_Get_Device_Interface_List_SizeW     = modcfgmgr32.NewProc("CM_Get_Device_Interface_List_SizeW")
	procCM_Locate_DevNodeW                     = modcfgmgr32.NewProc("CM_Locate_DevNodeW")
	procCimAddFsToMergedImage                  = modcimfs.NewProc("CimAddFsToMergedImage")
	procCimAddFsToMergedImage2                 = modcimfs.NewProc("CimAddFsToMergedImage2")
	procCimCloseImage                          = modcimfs.NewProc("CimCloseImage")
	procCimCloseStream                         = modcimfs.NewProc("CimCloseStream")
	procCimCommitImage                         = modcimfs.NewProc("CimCommitImage")
	procCimCreateAlternateStream               = modcimfs.NewProc("CimCreateAlternateStream")
	procCimCreateFile                          = modcimfs.NewProc("CimCreateFile")
	procCimCreateHardLink                      = modcimfs.NewProc("CimCreateHardLink")
	procCimCreateImage                         = modcimfs.NewProc("CimCreateImage")
	procCimCreateImage2                        = modcimfs.NewProc("CimCreateImage2")
	procCimCreateMergeLink                     = modcimfs.NewProc("CimCreateMergeLink")
	procCimDeletePath                          = modcimfs.NewProc("CimDeletePath")
	procCimDismountImage                       = modcimfs.NewProc("CimDismountImage")
	procCimGetVerificationInformation          = modcimfs.NewProc("CimGetVerificationInformation")
	procCimMergeMountImage                     = modcimfs.NewProc("CimMergeMountImage")
	procCimMountImage                          = modcimfs.NewProc("CimMountImage")
	procCimMountVerifiedImage                  = modcimfs.NewProc("CimMountVerifiedImage")
	procCimSealImage                           = modcimfs.NewProc("CimSealImage")
	procCimTombstoneFile                       = modcimfs.NewProc("CimTombstoneFile")
	procCimWriteStream                         = modcimfs.NewProc("CimWriteStream")
	procSetJobCompartmentId                    = modiphlpapi.NewProc("SetJobCompartmentId")
	procClosePseudoConsole                     = modkernel32.NewProc("ClosePseudoConsole")
	procCopyFileW                              = modkernel32.NewProc("CopyFileW")
	procCreatePseudoConsole                    = modkernel32.NewProc("CreatePseudoConsole")
	procCreateRemoteThread                     = modkernel32.NewProc("CreateRemoteThread")
	procGetActiveProcessorCount                = modkernel32.NewProc("GetActiveProcessorCount")
	procIsProcessInJob                         = modkernel32.NewProc("IsProcessInJob")
	procLocalAlloc                             = modkernel32.NewProc("LocalAlloc")
	procLocalFree                              = modkernel32.NewProc("LocalFree")
	procOpenJobObjectW                         = modkernel32.NewProc("OpenJobObjectW")
	procQueryInformationJobObject              = modkernel32.NewProc("QueryInformationJobObject")
	procQueryIoRateControlInformationJobObject = modkernel32.NewProc("QueryIoRateControlInformationJobObject")
	procResizePseudoConsole                    = modkernel32.NewProc("ResizePseudoConsole")
	procSearchPathW                            = modkernel32.NewProc("SearchPathW")
	procSetIoRateControlInformationJobObject   = modkernel32.NewProc("SetIoRateControlInformationJobObject")
	procNetLocalGroupAddMembers                = modnetapi32.NewProc("NetLocalGroupAddMembers")
	procNetLocalGroupGetInfo                   = modnetapi32.NewProc("NetLocalGroupGetInfo")
	procNetUserAdd                             = modnetapi32.NewProc("NetUserAdd")
	procNetUserDel                             = modnetapi32.NewProc("NetUserDel")
	procNtCreateFile                           = modntdll.NewProc("NtCreateFile")
	procNtCreateJobObject                      = modntdll.NewProc("NtCreateJobObject")
	procNtOpenDirectoryObject                  = modntdll.NewProc("NtOpenDirectoryObject")
	procNtOpenJobObject                        = modntdll.NewProc("NtOpenJobObject")
	procNtQueryDirectoryObject                 = modntdll.NewProc("NtQueryDirectoryObject")
	procNtQueryInformationProcess              = modntdll.NewProc("NtQueryInformationProcess")
	procNtQuerySystemInformation               = modntdll.NewProc("NtQuerySystemInformation")
	procNtSetInformationFile                   = modntdll.NewProc("NtSetInformationFile")
	procRtlNtStatusToDosError                  = modntdll.NewProc("RtlNtStatusToDosError")
	procORCloseHive                            = modoffreg.NewProc("ORCloseHive")
	procORCloseKey                             = modoffreg.NewProc("ORCloseKey")
	procORCreateHive                           = modoffreg.NewProc("ORCreateHive")
	procORCreateKey                            = modoffreg.NewProc("ORCreateKey")
	procORDeleteKey                            = modoffreg.NewProc("ORDeleteKey")
	procORGetValue                             = modoffreg.NewProc("ORGetValue")
	procORMergeHives                           = modoffreg.NewProc("ORMergeHives")
	procOROpenHive                             = modoffreg.NewProc("OROpenHive")
	procOROpenKey                              = modoffreg.NewProc("OROpenKey")
	procORSaveHive                             = modoffreg.NewProc("ORSaveHive")
	procORSetValue                             = modoffreg.NewProc("ORSetValue")
)

func LogonUser(username *uint16, domain *uint16, password *uint16, logonType uint32, logonProvider uint32, token *windows.Token) (err error) {
	r1, _, e1 := syscall.SyscallN(procLogonUserW.Addr(), uintptr(unsafe.Pointer(username)), uintptr(unsafe.Pointer(domain)), uintptr(unsafe.Pointer(password)), uintptr(logonType), uintptr(logonProvider), uintptr(unsafe.Pointer(token)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func BfSetupFilter(jobHandle windows.Handle, flags uint32, virtRootPath *uint16, virtTargetPath *uint16, virtExceptions **uint16, virtExceptionPathCount uint32) (hr error) {
	hr = procBfSetupFilter.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procBfSetupFilter.Addr(), uintptr(jobHandle), uintptr(flags), uintptr(unsafe.Pointer(virtRootPath)), uintptr(unsafe.Pointer(virtTargetPath)), uintptr(unsafe.Pointer(virtExceptions)), uintptr(virtExceptionPathCount))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CMGetDevNodeProperty(dnDevInst uint32, propertyKey *DevPropKey, propertyType *uint32, propertyBuffer *uint16, propertyBufferSize *uint32, uFlags uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procCM_Get_DevNode_PropertyW.Addr(), uintptr(dnDevInst), uintptr(unsafe.Pointer(propertyKey)), uintptr(unsafe.Pointer(propertyType)), uintptr(unsafe.Pointer(propertyBuffer)), uintptr(unsafe.Pointer(propertyBufferSize)), uintptr(uFlags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CMGetDeviceIDList(pszFilter *byte, buffer *byte, bufferLen uint32, uFlags uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procCM_Get_Device_ID_ListA.Addr(), uintptr(unsafe.Pointer(pszFilter)), uintptr(unsafe.Pointer(buffer)), uintptr(bufferLen), uintptr(uFlags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CMGetDeviceIDListSize(pulLen *uint32, pszFilter *byte, uFlags uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procCM_Get_Device_ID_List_SizeA.Addr(), uintptr(unsafe.Pointer(pulLen)), uintptr(unsafe.Pointer(pszFilter)), uintptr(uFlags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CMGetDeviceInterfaceList(classGUID *g, deviceID *uint16, buffer *uint16, bufLen uint32, ulFlags uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procCM_Get_Device_Interface_ListW.Addr(), uintptr(unsafe.Pointer(classGUID)), uintptr(unsafe.Pointer(deviceID)), uintptr(unsafe.Pointer(buffer)), uintptr(bufLen), uintptr(ulFlags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CMGetDeviceInterfaceListSize(listlen *uint32, classGUID *g, deviceID *uint16, ulFlags uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procCM_Get_Device_Interface_List_SizeW.Addr(), uintptr(unsafe.Pointer(listlen)), uintptr(unsafe.Pointer(classGUID)), uintptr(unsafe.Pointer(deviceID)), uintptr(ulFlags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CMLocateDevNode(pdnDevInst *uint32, pDeviceID string, uFlags uint32) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(pDeviceID)
	if hr != nil {
		return
	}
	return _CMLocateDevNode(pdnDevInst, _p0, uFlags)
}

func _CMLocateDevNode(pdnDevInst *uint32, pDeviceID *uint16, uFlags uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procCM_Locate_DevNodeW.Addr(), uintptr(unsafe.Pointer(pdnDevInst)), uintptr(unsafe.Pointer(pDeviceID)), uintptr(uFlags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimAddFsToMergedImage(cimFSHandle FsHandle, path string) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(path)
	if hr != nil {
		return
	}
	return _CimAddFsToMergedImage(cimFSHandle, _p0)
}

func _CimAddFsToMergedImage(cimFSHandle FsHandle, path *uint16) (hr error) {
	hr = procCimAddFsToMergedImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimAddFsToMergedImage.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(path)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimAddFsToMergedImage2(cimFSHandle FsHandle, path string, flags uint32) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(path)
	if hr != nil {
		return
	}
	return _CimAddFsToMergedImage2(cimFSHandle, _p0, flags)
}

func _CimAddFsToMergedImage2(cimFSHandle FsHandle, path *uint16, flags uint32) (hr error) {
	hr = procCimAddFsToMergedImage2.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimAddFsToMergedImage2.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(path)), uintptr(flags))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCloseImage(cimFSHandle FsHandle) (err error) {
	err = procCimCloseImage.Find()
	if err != nil {
		return
	}
	syscall.SyscallN(procCimCloseImage.Addr(), uintptr(cimFSHandle))
	return
}

func CimCloseStream(cimStreamHandle StreamHandle) (hr error) {
	hr = procCimCloseStream.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCloseStream.Addr(), uintptr(cimStreamHandle))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCommitImage(cimFSHandle FsHandle) (hr error) {
	hr = procCimCommitImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCommitImage.Addr(), uintptr(cimFSHandle))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCreateAlternateStream(cimFSHandle FsHandle, path string, size uint64, cimStreamHandle *StreamHandle) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(path)
	if hr != nil {
		return
	}
	return _CimCreateAlternateStream(cimFSHandle, _p0, size, cimStreamHandle)
}

func _CimCreateAlternateStream(cimFSHandle FsHandle, path *uint16, size uint64, cimStreamHandle *StreamHandle) (hr error) {
	hr = procCimCreateAlternateStream.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCreateAlternateStream.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(path)), uintptr(size), uintptr(unsafe.Pointer(cimStreamHandle)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCreateFile(cimFSHandle FsHandle, path string, file *CimFsFileMetadata, cimStreamHandle *StreamHandle) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(path)
	if hr != nil {
		return
	}
	return _CimCreateFile(cimFSHandle, _p0, file, cimStreamHandle)
}

func _CimCreateFile(cimFSHandle FsHandle, path *uint16, file *CimFsFileMetadata, cimStreamHandle *StreamHandle) (hr error) {
	hr = procCimCreateFile.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCreateFile.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(cimStreamHandle)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCreateHardLink(cimFSHandle FsHandle, newPath string, oldPath string) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(newPath)
	if hr != nil {
		return
	}
	var _p1 *uint16
	_p1, hr = syscall.UTF16PtrFromString(oldPath)
	if hr != nil {
		return
	}
	return _CimCreateHardLink(cimFSHandle, _p0, _p1)
}

func _CimCreateHardLink(cimFSHandle FsHandle, newPath *uint16, oldPath *uint16) (hr error) {
	hr = procCimCreateHardLink.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCreateHardLink.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(newPath)), uintptr(unsafe.Pointer(oldPath)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCreateImage(imagePath string, oldFSName *uint16, newFSName *uint16, cimFSHandle *FsHandle) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(imagePath)
	if hr != nil {
		return
	}
	return _CimCreateImage(_p0, oldFSName, newFSName, cimFSHandle)
}

func _CimCreateImage(imagePath *uint16, oldFSName *uint16, newFSName *uint16, cimFSHandle *FsHandle) (hr error) {
	hr = procCimCreateImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCreateImage.Addr(), uintptr(unsafe.Pointer(imagePath)), uintptr(unsafe.Pointer(oldFSName)), uintptr(unsafe.Pointer(newFSName)), uintptr(unsafe.Pointer(cimFSHandle)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCreateImage2(imagePath string, flags uint32, oldFSName *uint16, newFSName *uint16, cimFSHandle *FsHandle) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(imagePath)
	if hr != nil {
		return
	}
	return _CimCreateImage2(_p0, flags, oldFSName, newFSName, cimFSHandle)
}

func _CimCreateImage2(imagePath *uint16, flags uint32, oldFSName *uint16, newFSName *uint16, cimFSHandle *FsHandle) (hr error) {
	hr = procCimCreateImage2.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCreateImage2.Addr(), uintptr(unsafe.Pointer(imagePath)), uintptr(flags), uintptr(unsafe.Pointer(oldFSName)), uintptr(unsafe.Pointer(newFSName)), uintptr(unsafe.Pointer(cimFSHandle)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimCreateMergeLink(cimFSHandle FsHandle, newPath string, oldPath string) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(newPath)
	if hr != nil {
		return
	}
	var _p1 *uint16
	_p1, hr = syscall.UTF16PtrFromString(oldPath)
	if hr != nil {
		return
	}
	return _CimCreateMergeLink(cimFSHandle, _p0, _p1)
}

func _CimCreateMergeLink(cimFSHandle FsHandle, newPath *uint16, oldPath *uint16) (hr error) {
	hr = procCimCreateMergeLink.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimCreateMergeLink.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(newPath)), uintptr(unsafe.Pointer(oldPath)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimDeletePath(cimFSHandle FsHandle, path string) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(path)
	if hr != nil {
		return
	}
	return _CimDeletePath(cimFSHandle, _p0)
}

func _CimDeletePath(cimFSHandle FsHandle, path *uint16) (hr error) {
	hr = procCimDeletePath.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimDeletePath.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(path)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimDismountImage(volumeID *g) (hr error) {
	hr = procCimDismountImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimDismountImage.Addr(), uintptr(unsafe.Pointer(volumeID)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimGetVerificationInformation(blockCimPath string, isSealed *uint32, hashSize *uint64, signatureSize *uint64, fixedHeaderSize *uint64, hash *byte, signature *byte) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(blockCimPath)
	if hr != nil {
		return
	}
	return _CimGetVerificationInformation(_p0, isSealed, hashSize, signatureSize, fixedHeaderSize, hash, signature)
}

func _CimGetVerificationInformation(blockCimPath *uint16, isSealed *uint32, hashSize *uint64, signatureSize *uint64, fixedHeaderSize *uint64, hash *byte, signature *byte) (hr error) {
	hr = procCimGetVerificationInformation.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimGetVerificationInformation.Addr(), uintptr(unsafe.Pointer(blockCimPath)), uintptr(unsafe.Pointer(isSealed)), uintptr(unsafe.Pointer(hashSize)), uintptr(unsafe.Pointer(signatureSize)), uintptr(unsafe.Pointer(fixedHeaderSize)), uintptr(unsafe.Pointer(hash)), uintptr(unsafe.Pointer(signature)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimMergeMountImage(numCimPaths uint32, backingImagePaths *CimFsImagePath, flags uint32, volumeID *g) (hr error) {
	hr = procCimMergeMountImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimMergeMountImage.Addr(), uintptr(numCimPaths), uintptr(unsafe.Pointer(backingImagePaths)), uintptr(flags), uintptr(unsafe.Pointer(volumeID)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimMountImage(imagePath string, fsName string, flags uint32, volumeID *g) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(imagePath)
	if hr != nil {
		return
	}
	var _p1 *uint16
	_p1, hr = syscall.UTF16PtrFromString(fsName)
	if hr != nil {
		return
	}
	return _CimMountImage(_p0, _p1, flags, volumeID)
}

func _CimMountImage(imagePath *uint16, fsName *uint16, flags uint32, volumeID *g) (hr error) {
	hr = procCimMountImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimMountImage.Addr(), uintptr(unsafe.Pointer(imagePath)), uintptr(unsafe.Pointer(fsName)), uintptr(flags), uintptr(unsafe.Pointer(volumeID)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimMountVerifiedImage(imagePath string, fsName string, flags uint32, volumeID *g, hashSize uint16, hash *byte) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(imagePath)
	if hr != nil {
		return
	}
	var _p1 *uint16
	_p1, hr = syscall.UTF16PtrFromString(fsName)
	if hr != nil {
		return
	}
	return _CimMountVerifiedImage(_p0, _p1, flags, volumeID, hashSize, hash)
}

func _CimMountVerifiedImage(imagePath *uint16, fsName *uint16, flags uint32, volumeID *g, hashSize uint16, hash *byte) (hr error) {
	hr = procCimMountVerifiedImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimMountVerifiedImage.Addr(), uintptr(unsafe.Pointer(imagePath)), uintptr(unsafe.Pointer(fsName)), uintptr(flags), uintptr(unsafe.Pointer(volumeID)), uintptr(hashSize), uintptr(unsafe.Pointer(hash)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimSealImage(blockCimPath string, hashSize *uint64, fixedHeaderSize *uint64, hash *byte) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(blockCimPath)
	if hr != nil {
		return
	}
	return _CimSealImage(_p0, hashSize, fixedHeaderSize, hash)
}

func _CimSealImage(blockCimPath *uint16, hashSize *uint64, fixedHeaderSize *uint64, hash *byte) (hr error) {
	hr = procCimSealImage.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimSealImage.Addr(), uintptr(unsafe.Pointer(blockCimPath)), uintptr(unsafe.Pointer(hashSize)), uintptr(unsafe.Pointer(fixedHeaderSize)), uintptr(unsafe.Pointer(hash)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimTombstoneFile(cimFSHandle FsHandle, path string) (hr error) {
	var _p0 *uint16
	_p0, hr = syscall.UTF16PtrFromString(path)
	if hr != nil {
		return
	}
	return _CimTombstoneFile(cimFSHandle, _p0)
}

func _CimTombstoneFile(cimFSHandle FsHandle, path *uint16) (hr error) {
	hr = procCimTombstoneFile.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimTombstoneFile.Addr(), uintptr(cimFSHandle), uintptr(unsafe.Pointer(path)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CimWriteStream(cimStreamHandle StreamHandle, buffer uintptr, bufferSize uint32) (hr error) {
	hr = procCimWriteStream.Find()
	if hr != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(procCimWriteStream.Addr(), uintptr(cimStreamHandle), uintptr(buffer), uintptr(bufferSize))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func SetJobCompartmentId(handle windows.Handle, compartmentId uint32) (win32Err error) {
	r0, _, _ := syscall.SyscallN(procSetJobCompartmentId.Addr(), uintptr(handle), uintptr(compartmentId))
	if r0 != 0 {
		win32Err = syscall.Errno(r0)
	}
	return
}

func ClosePseudoConsole(hpc windows.Handle) {
	syscall.SyscallN(procClosePseudoConsole.Addr(), uintptr(hpc))
	return
}

func CopyFileW(existingFileName *uint16, newFileName *uint16, failIfExists int32) (err error) {
	r1, _, e1 := syscall.SyscallN(procCopyFileW.Addr(), uintptr(unsafe.Pointer(existingFileName)), uintptr(unsafe.Pointer(newFileName)), uintptr(failIfExists))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func createPseudoConsole(size uint32, hInput windows.Handle, hOutput windows.Handle, dwFlags uint32, hpcon *windows.Handle) (hr error) {
	r0, _, _ := syscall.SyscallN(procCreatePseudoConsole.Addr(), uintptr(size), uintptr(hInput), uintptr(hOutput), uintptr(dwFlags), uintptr(unsafe.Pointer(hpcon)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func CreateRemoteThread(process windows.Handle, sa *windows.SecurityAttributes, stackSize uint32, startAddr uintptr, parameter uintptr, creationFlags uint32, threadID *uint32) (handle windows.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(procCreateRemoteThread.Addr(), uintptr(process), uintptr(unsafe.Pointer(sa)), uintptr(stackSize), uintptr(startAddr), uintptr(parameter), uintptr(creationFlags), uintptr(unsafe.Pointer(threadID)))
	handle = windows.Handle(r0)
	if handle == 0 {
		err = errnoErr(e1)
	}
	return
}

func GetActiveProcessorCount(groupNumber uint16) (amount uint32) {
	r0, _, _ := syscall.SyscallN(procGetActiveProcessorCount.Addr(), uintptr(groupNumber))
	amount = uint32(r0)
	return
}

func IsProcessInJob(procHandle windows.Handle, jobHandle windows.Handle, result *int32) (err error) {
	r1, _, e1 := syscall.SyscallN(procIsProcessInJob.Addr(), uintptr(procHandle), uintptr(jobHandle), uintptr(unsafe.Pointer(result)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func LocalAlloc(flags uint32, size int) (ptr uintptr) {
	r0, _, _ := syscall.SyscallN(procLocalAlloc.Addr(), uintptr(flags), uintptr(size))
	ptr = uintptr(r0)
	return
}

func LocalFree(ptr uintptr) {
	syscall.SyscallN(procLocalFree.Addr(), uintptr(ptr))
	return
}

func OpenJobObject(desiredAccess uint32, inheritHandle bool, lpName *uint16) (handle windows.Handle, err error) {
	var _p0 uint32
	if inheritHandle {
		_p0 = 1
	}
	r0, _, e1 := syscall.SyscallN(procOpenJobObjectW.Addr(), uintptr(desiredAccess), uintptr(_p0), uintptr(unsafe.Pointer(lpName)))
	handle = windows.Handle(r0)
	if handle == 0 {
		err = errnoErr(e1)
	}
	return
}

func QueryInformationJobObject(jobHandle windows.Handle, infoClass uint32, jobObjectInfo unsafe.Pointer, jobObjectInformationLength uint32, lpReturnLength *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procQueryInformationJobObject.Addr(), uintptr(jobHandle), uintptr(infoClass), uintptr(jobObjectInfo), uintptr(jobObjectInformationLength), uintptr(unsafe.Pointer(lpReturnLength)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func QueryIoRateControlInformationJobObject(jobHandle windows.Handle, volumeName *uint16, ioRateControlInfo **JOBOBJECT_IO_RATE_CONTROL_INFORMATION, infoBlockCount *uint32) (ret uint32, err error) {
	r0, _, e1 := syscall.SyscallN(procQueryIoRateControlInformationJobObject.Addr(), uintptr(jobHandle), uintptr(unsafe.Pointer(volumeName)), uintptr(unsafe.Pointer(ioRateControlInfo)), uintptr(unsafe.Pointer(infoBlockCount)))
	ret = uint32(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func resizePseudoConsole(hPc windows.Handle, size uint32) (hr error) {
	r0, _, _ := syscall.SyscallN(procResizePseudoConsole.Addr(), uintptr(hPc), uintptr(size))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

func SearchPath(lpPath *uint16, lpFileName *uint16, lpExtension *uint16, nBufferLength uint32, lpBuffer *uint16, lpFilePath *uint16) (size uint32, err error) {
	r0, _, e1 := syscall.SyscallN(procSearchPathW.Addr(), uintptr(unsafe.Pointer(lpPath)), uintptr(unsafe.Pointer(lpFileName)), uintptr(unsafe.Pointer(lpExtension)), uintptr(nBufferLength), uintptr(unsafe.Pointer(lpBuffer)), uintptr(unsafe.Pointer(lpFilePath)))
	size = uint32(r0)
	if size == 0 {
		err = errnoErr(e1)
	}
	return
}

func SetIoRateControlInformationJobObject(jobHandle windows.Handle, ioRateControlInfo *JOBOBJECT_IO_RATE_CONTROL_INFORMATION) (ret uint32, err error) {
	r0, _, e1 := syscall.SyscallN(procSetIoRateControlInformationJobObject.Addr(), uintptr(jobHandle), uintptr(unsafe.Pointer(ioRateControlInfo)))
	ret = uint32(r0)
	if ret == 0 {
		err = errnoErr(e1)
	}
	return
}

func netLocalGroupAddMembers(serverName *uint16, groupName *uint16, level uint32, buf *byte, totalEntries uint32) (status error) {
	r0, _, _ := syscall.SyscallN(procNetLocalGroupAddMembers.Addr(), uintptr(unsafe.Pointer(serverName)), uintptr(unsafe.Pointer(groupName)), uintptr(level), uintptr(unsafe.Pointer(buf)), uintptr(totalEntries))
	if r0 != 0 {
		status = syscall.Errno(r0)
	}
	return
}

func netLocalGroupGetInfo(serverName *uint16, groupName *uint16, level uint32, bufptr **byte) (status error) {
	r0, _, _ := syscall.SyscallN(procNetLocalGroupGetInfo.Addr(), uintptr(unsafe.Pointer(serverName)), uintptr(unsafe.Pointer(groupName)), uintptr(level), uintptr(unsafe.Pointer(bufptr)))
	if r0 != 0 {
		status = syscall.Errno(r0)
	}
	return
}

func netUserAdd(serverName *uint16, level uint32, buf *byte, parm_err *uint32) (status error) {
	r0, _, _ := syscall.SyscallN(procNetUserAdd.Addr(), uintptr(unsafe.Pointer(serverName)), uintptr(level), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(parm_err)))
	if r0 != 0 {
		status = syscall.Errno(r0)
	}
	return
}

func netUserDel(serverName *uint16, username *uint16) (status error) {
	r0, _, _ := syscall.SyscallN(procNetUserDel.Addr(), uintptr(unsafe.Pointer(serverName)), uintptr(unsafe.Pointer(username)))
	if r0 != 0 {
		status = syscall.Errno(r0)
	}
	return
}

func NtCreateFile(handle *uintptr, accessMask uint32, oa *ObjectAttributes, iosb *IOStatusBlock, allocationSize *uint64, fileAttributes uint32, shareAccess uint32, createDisposition uint32, createOptions uint32, eaBuffer *byte, eaLength uint32) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtCreateFile.Addr(), uintptr(unsafe.Pointer(handle)), uintptr(accessMask), uintptr(unsafe.Pointer(oa)), uintptr(unsafe.Pointer(iosb)), uintptr(unsafe.Pointer(allocationSize)), uintptr(fileAttributes), uintptr(shareAccess), uintptr(createDisposition), uintptr(createOptions), uintptr(unsafe.Pointer(eaBuffer)), uintptr(eaLength))
	status = uint32(r0)
	return
}

func NtCreateJobObject(jobHandle *windows.Handle, desiredAccess uint32, objAttributes *ObjectAttributes) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtCreateJobObject.Addr(), uintptr(unsafe.Pointer(jobHandle)), uintptr(desiredAccess), uintptr(unsafe.Pointer(objAttributes)))
	status = uint32(r0)
	return
}

func NtOpenDirectoryObject(handle *uintptr, accessMask uint32, oa *ObjectAttributes) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtOpenDirectoryObject.Addr(), uintptr(unsafe.Pointer(handle)), uintptr(accessMask), uintptr(unsafe.Pointer(oa)))
	status = uint32(r0)
	return
}

func NtOpenJobObject(jobHandle *windows.Handle, desiredAccess uint32, objAttributes *ObjectAttributes) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtOpenJobObject.Addr(), uintptr(unsafe.Pointer(jobHandle)), uintptr(desiredAccess), uintptr(unsafe.Pointer(objAttributes)))
	status = uint32(r0)
	return
}

func NtQueryDirectoryObject(handle uintptr, buffer *byte, length uint32, singleEntry bool, restartScan bool, context *uint32, returnLength *uint32) (status uint32) {
	var _p0 uint32
	if singleEntry {
		_p0 = 1
	}
	var _p1 uint32
	if restartScan {
		_p1 = 1
	}
	r0, _, _ := syscall.SyscallN(procNtQueryDirectoryObject.Addr(), uintptr(handle), uintptr(unsafe.Pointer(buffer)), uintptr(length), uintptr(_p0), uintptr(_p1), uintptr(unsafe.Pointer(context)), uintptr(unsafe.Pointer(returnLength)))
	status = uint32(r0)
	return
}

func NtQueryInformationProcess(processHandle windows.Handle, processInfoClass uint32, processInfo unsafe.Pointer, processInfoLength uint32, returnLength *uint32) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtQueryInformationProcess.Addr(), uintptr(processHandle), uintptr(processInfoClass), uintptr(processInfo), uintptr(processInfoLength), uintptr(unsafe.Pointer(returnLength)))
	status = uint32(r0)
	return
}

func NtQuerySystemInformation(systemInfoClass int, systemInformation unsafe.Pointer, systemInfoLength uint32, returnLength *uint32) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtQuerySystemInformation.Addr(), uintptr(systemInfoClass), uintptr(systemInformation), uintptr(systemInfoLength), uintptr(unsafe.Pointer(returnLength)))
	status = uint32(r0)
	return
}

func NtSetInformationFile(handle uintptr, iosb *IOStatusBlock, information uintptr, length uint32, class uint32) (status uint32) {
	r0, _, _ := syscall.SyscallN(procNtSetInformationFile.Addr(), uintptr(handle), uintptr(unsafe.Pointer(iosb)), uintptr(information), uintptr(length), uintptr(class))
	status = uint32(r0)
	return
}

func RtlNtStatusToDosError(status uint32) (winerr error) {
	r0, _, _ := syscall.SyscallN(procRtlNtStatusToDosError.Addr(), uintptr(status))
	if r0 != 0 {
		winerr = syscall.Errno(r0)
	}
	return
}

func ORCloseHive(handle ORHKey) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORCloseHive.Addr(), uintptr(handle))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORCloseKey(handle ORHKey) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORCloseKey.Addr(), uintptr(handle))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORCreateHive(key *ORHKey) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORCreateHive.Addr(), uintptr(unsafe.Pointer(key)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORCreateKey(handle ORHKey, subKey string, class uintptr, options uint32, securityDescriptor uintptr, result *ORHKey, disposition *uint32) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(subKey)
	if win32err != nil {
		return
	}
	return _ORCreateKey(handle, _p0, class, options, securityDescriptor, result, disposition)
}

func _ORCreateKey(handle ORHKey, subKey *uint16, class uintptr, options uint32, securityDescriptor uintptr, result *ORHKey, disposition *uint32) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORCreateKey.Addr(), uintptr(handle), uintptr(unsafe.Pointer(subKey)), uintptr(class), uintptr(options), uintptr(securityDescriptor), uintptr(unsafe.Pointer(result)), uintptr(unsafe.Pointer(disposition)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORDeleteKey(handle ORHKey, subKey string) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(subKey)
	if win32err != nil {
		return
	}
	return _ORDeleteKey(handle, _p0)
}

func _ORDeleteKey(handle ORHKey, subKey *uint16) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORDeleteKey.Addr(), uintptr(handle), uintptr(unsafe.Pointer(subKey)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORGetValue(handle ORHKey, subKey string, value string, valueType *uint32, data *byte, dataLen *uint32) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(subKey)
	if win32err != nil {
		return
	}
	var _p1 *uint16
	_p1, win32err = syscall.UTF16PtrFromString(value)
	if win32err != nil {
		return
	}
	return _ORGetValue(handle, _p0, _p1, valueType, data, dataLen)
}

func _ORGetValue(handle ORHKey, subKey *uint16, value *uint16, valueType *uint32, data *byte, dataLen *uint32) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORGetValue.Addr(), uintptr(handle), uintptr(unsafe.Pointer(subKey)), uintptr(unsafe.Pointer(value)), uintptr(unsafe.Pointer(valueType)), uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(dataLen)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORMergeHives(hiveHandles []ORHKey, result *ORHKey) (win32err error) {
	var _p0 *ORHKey
	if len(hiveHandles) > 0 {
		_p0 = &hiveHandles[0]
	}
	r0, _, _ := syscall.SyscallN(procORMergeHives.Addr(), uintptr(unsafe.Pointer(_p0)), uintptr(len(hiveHandles)), uintptr(unsafe.Pointer(result)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func OROpenHive(hivePath string, result *ORHKey) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(hivePath)
	if win32err != nil {
		return
	}
	return _OROpenHive(_p0, result)
}

func _OROpenHive(hivePath *uint16, result *ORHKey) (win32err error) {
	r0, _, _ := syscall.SyscallN(procOROpenHive.Addr(), uintptr(unsafe.Pointer(hivePath)), uintptr(unsafe.Pointer(result)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func OROpenKey(handle ORHKey, subKey string, result *ORHKey) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(subKey)
	if win32err != nil {
		return
	}
	return _OROpenKey(handle, _p0, result)
}

func _OROpenKey(handle ORHKey, subKey *uint16, result *ORHKey) (win32err error) {
	r0, _, _ := syscall.SyscallN(procOROpenKey.Addr(), uintptr(handle), uintptr(unsafe.Pointer(subKey)), uintptr(unsafe.Pointer(result)))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORSaveHive(handle ORHKey, hivePath string, osMajorVersion uint32, osMinorVersion uint32) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(hivePath)
	if win32err != nil {
		return
	}
	return _ORSaveHive(handle, _p0, osMajorVersion, osMinorVersion)
}

func _ORSaveHive(handle ORHKey, hivePath *uint16, osMajorVersion uint32, osMinorVersion uint32) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORSaveHive.Addr(), uintptr(handle), uintptr(unsafe.Pointer(hivePath)), uintptr(osMajorVersion), uintptr(osMinorVersion))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}

func ORSetValue(handle ORHKey, valueName string, valueType uint32, data *byte, dataLen uint32) (win32err error) {
	var _p0 *uint16
	_p0, win32err = syscall.UTF16PtrFromString(valueName)
	if win32err != nil {
		return
	}
	return _ORSetValue(handle, _p0, valueType, data, dataLen)
}

func _ORSetValue(handle ORHKey, valueName *uint16, valueType uint32, data *byte, dataLen uint32) (win32err error) {
	r0, _, _ := syscall.SyscallN(procORSetValue.Addr(), uintptr(handle), uintptr(unsafe.Pointer(valueName)), uintptr(valueType), uintptr(unsafe.Pointer(data)), uintptr(dataLen))
	if r0 != 0 {
		win32err = syscall.Errno(r0)
	}
	return
}
