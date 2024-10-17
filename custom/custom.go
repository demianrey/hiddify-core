package main

/*
#include "stdint.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"github.com/hiddify/hiddify-core/bridge"
	"github.com/hiddify/hiddify-core/config"

	hcore "github.com/hiddify/hiddify-core/v2/hcore"

	"github.com/sagernet/sing-box/log"
)

//export setupOnce
func setupOnce(api unsafe.Pointer) {
	bridge.InitializeDartApi(api)
}

//export setup
func setup(baseDir *C.char, workingDir *C.char, tempDir *C.char, mode C.int, listen *C.char, secret *C.char, statusPort C.longlong, debug bool) (CErr *C.char) {
	// err := hcore.Setup(C.GoString(baseDir), C.GoString(workingDir), C.GoString(tempDir), int64(statusPort), debug)
	err := hcore.Setup(hcore.SetupParameters{
		BasePath:          C.GoString(baseDir),
		WorkingDir:        C.GoString(workingDir),
		TempDir:           C.GoString(tempDir),
		FlutterStatusPort: int64(statusPort),
		Debug:             debug,
		Mode:              hcore.SetupMode(mode),
		Listen:            C.GoString(listen),
		Secret:            C.GoString(secret),
	})
	return emptyOrErrorC(err)
}

//export parse
func parse(path *C.char, tempPath *C.char, debug bool) (CErr *C.char) {
	res, err := hcore.Parse(&hcore.ParseRequest{
		ConfigPath: C.GoString(path),
		TempPath:   C.GoString(tempPath),
	})
	if err != nil {
		log.Error(err.Error())
		return C.CString(err.Error())
	}

	err = os.WriteFile(C.GoString(path), []byte(res.Content), 0o644)
	return emptyOrErrorC(err)
}

//export changeHiddifyOptions
func changeHiddifyOptions(HiddifyOptionsJson *C.char) (CErr *C.char) {
	_, err := hcore.ChangeHiddifySettings(&hcore.ChangeHiddifySettingsRequest{
		HiddifySettingsJson: C.GoString(HiddifyOptionsJson),
	})
	return emptyOrErrorC(err)
}

//export generateConfig
func generateConfig(path *C.char) (res *C.char) {
	conf, err := hcore.GenerateConfig(&hcore.GenerateConfigRequest{
		Path: C.GoString(path),
	})
	if err != nil {
		return emptyOrErrorC(err)
	}
	fmt.Printf("Config: %+v\n", conf)
	fmt.Printf("ConfigContent: %+v\n", conf.ConfigContent)
	return C.CString(conf.ConfigContent)
}

//export start
func start(configPath *C.char, disableMemoryLimit bool) (CErr *C.char) {
	_, err := hcore.Start(&hcore.StartRequest{
		ConfigPath:             C.GoString(configPath),
		EnableOldCommandServer: true,
		DisableMemoryLimit:     disableMemoryLimit,
	})
	return emptyOrErrorC(err)
}

//export stop
func stop() (CErr *C.char) {
	_, err := hcore.Stop()
	return emptyOrErrorC(err)
}

//export restart
func restart(configPath *C.char, disableMemoryLimit bool) (CErr *C.char) {
	_, err := hcore.Restart(&hcore.StartRequest{
		ConfigPath:             C.GoString(configPath),
		EnableOldCommandServer: true,
		DisableMemoryLimit:     disableMemoryLimit,
	})
	return emptyOrErrorC(err)
}

//export startCommandClient
func startCommandClient(command C.int, port C.longlong) *C.char {
	err := hcore.StartCommand(int32(command), int64(port))
	return emptyOrErrorC(err)
}

//export stopCommandClient
func stopCommandClient(command C.int) *C.char {
	err := hcore.StopCommand(int32(command))
	return emptyOrErrorC(err)
}

//export selectOutbound
func selectOutbound(groupTag *C.char, outboundTag *C.char) (CErr *C.char) {
	_, err := hcore.SelectOutbound(&hcore.SelectOutboundRequest{
		GroupTag:    C.GoString(groupTag),
		OutboundTag: C.GoString(outboundTag),
	})

	return emptyOrErrorC(err)
}

//export urlTest
func urlTest(groupTag *C.char) (CErr *C.char) {
	_, err := hcore.UrlTest(&hcore.UrlTestRequest{
		GroupTag: C.GoString(groupTag),
	})

	return emptyOrErrorC(err)
}

func emptyOrErrorC(err error) *C.char {
	if err == nil {
		return C.CString("")
	}
	log.Error(err.Error())
	return C.CString(err.Error())
}

//export generateWarpConfig
func generateWarpConfig(licenseKey *C.char, accountId *C.char, accessToken *C.char) (CResp *C.char) {
	res, err := hcore.GenerateWarpConfig(&hcore.GenerateWarpConfigRequest{
		LicenseKey:  C.GoString(licenseKey),
		AccountId:   C.GoString(accountId),
		AccessToken: C.GoString(accessToken),
	})
	if err != nil {
		return C.CString(fmt.Sprint("error: ", err.Error()))
	}
	warpAccount := config.WarpAccount{
		AccountID:   res.Account.AccountId,
		AccessToken: res.Account.AccessToken,
	}
	warpConfig := config.WarpWireguardConfig{
		PrivateKey:       res.Config.PrivateKey,
		LocalAddressIPv4: res.Config.LocalAddressIpv4,
		LocalAddressIPv6: res.Config.LocalAddressIpv6,
		PeerPublicKey:    res.Config.PeerPublicKey,
		ClientID:         res.Config.ClientId,
	}
	log := res.Log
	response := &config.WarpGenerationResponse{
		WarpAccount: warpAccount,
		Log:         log,
		Config:      warpConfig,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(responseJson))
}

func main() {}

//export GetServerPublicKey
func GetServerPublicKey() []byte {
	return hcore.GetGrpcServerPublicKey()
}

//export AddGrpcClientPublicKey
func AddGrpcClientPublicKey(clientPublicKey []byte) error {
	return hcore.AddGrpcClientPublicKey(clientPublicKey)
}

//export close
func close(mode C.int) {
	hcore.Close(hcore.SetupMode(int32(mode)))
}
