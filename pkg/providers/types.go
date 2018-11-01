package providers

type ResourceAction string

const (
	ResourceCreated ResourceAction = "create"
	ResourceUpdated ResourceAction = "update"
	ResourceDeleted ResourceAction = "delete"
)

type ResourceReference struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
}

type ResourceEvent struct {
	Action   ResourceAction    `json:"action"`
	Resource ResourceReference `json:"resource"`
}
