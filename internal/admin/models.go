package admin

type info struct {
	InstanceUUID string `json:"instance_uuid"`
	ServiceName  string `json:"service_name"`
	AppName      string `json:"app_name"`
	GitHash      string `json:"git_hash"`
	Version      string `json:"version"`
	BuildAt      string `json:"build_at"`
	StartTime    string `json:"start_time"`
	UpTime       string `json:"up_time"`
}
