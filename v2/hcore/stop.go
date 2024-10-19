package hcore

import (
	"context"
	"fmt"

	common "github.com/hiddify/hiddify-core/v2/common"
	"github.com/hiddify/hiddify-core/v2/config"
)

func (s *CoreService) Stop(ctx context.Context, empty *common.Empty) (*CoreInfoResponse, error) {
	return Stop()
}

func Stop() (coreResponse *CoreInfoResponse, err error) {
	defer config.DeferPanicToError("stop", func(recovered_err error) {
		coreResponse, err = errorWrapper(MessageType_UNEXPECTED_ERROR, recovered_err)
	})

	if CoreState != CoreStates_STARTED {
		return errorWrapper(MessageType_INSTANCE_NOT_STARTED, fmt.Errorf("instance not started"))
	}
	if Box == nil {
		return errorWrapper(MessageType_INSTANCE_NOT_FOUND, fmt.Errorf("instance not found"))
	}
	SetCoreStatus(CoreStates_STOPPING, MessageType_EMPTY, "")

	err = Box.Close()
	if err != nil {
		return errorWrapper(MessageType_UNEXPECTED_ERROR, err)
	}
	Box = nil
	return SetCoreStatus(CoreStates_STOPPED, MessageType_EMPTY, ""), nil
}
