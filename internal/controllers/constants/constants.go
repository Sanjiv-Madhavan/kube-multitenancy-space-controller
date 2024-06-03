package constants

import "time"

const (
	SpaceFinalizer string = "github.sanjivmadhavan.io/kube-multitenancy-space-finalizer"
	RequeueAfter          = time.Minute * 3

	SpaceConditionReady    string = "Ready"
	SpaceConditionCreating string = "Creating"
	SpaceConditionFailed   string = "Failed"

	SpaceSyncSuccessReason string = "SpaceSyncedSuccessfully"
	SpaceCreatingReason    string = "SpaceCreating"
	SpaceFailedReason      string = "SpaceSyncFailed"

	SpaceSyncSuccessMessage string = "Space synced successfully"
	SpaceSyncFailMessage    string = "Space failed to sync"
	SpaceCreatingMessage    string = "Creating Space in progress"
)
