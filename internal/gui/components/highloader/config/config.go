package config

// PreloadConfig struct for presetting form from config.
type PreloadConfig struct {
	Host  string  `json:"host"`
	Proto []Proto `json:"proto"`
}

// Proto struct with info one proto file.
type Proto struct {
	FilePath string    `json:"file_path"`
	Requests []Request `json:"requests"`
}

// Request struct with info one request for proto.
type Request struct {
	LogPath         string     `json:"log_path"`
	MetricsPath     string     `json:"metrics_path"`
	Message         string     `json:"message"`
	StopAfter       Time       `json:"stop_after"`
	RequestDeadline Time       `json:"request_deadline"`
	Service         string     `json:"service"`
	Method          string     `json:"method"`
	RPS             string     `json:"rps"`
	Metadata        []MetaData `json:"metadata"`
}

// MetaData struct with metadata for request.
type MetaData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Time type time for requests settings.
type Time struct {
	Duration string `json:"duration"`
	Type     string `json:"type"`
}
