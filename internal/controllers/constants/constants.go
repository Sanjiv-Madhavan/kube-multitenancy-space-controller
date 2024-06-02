package constants

const (
	SpaceFinalizer string = "github.sanjivmadhavan.io/kube-multitenancy-space-finalizer"

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
