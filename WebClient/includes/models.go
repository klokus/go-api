package includes

// StatusResponse is
type StatusResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// AttackRequest is
type AttackRequest struct {
	Host   string `json:"host"`
	Time   int    `json:"time"`
	Port   int    `json:"port"`
	Method string `json:"method"`
}

// ConfigStruct is our struct for the JSON config
type ConfigStruct struct {
	Methods []Method `json:"methods"`
	MaxTime int      `json:"max_time"`
	Servers []Server `json:"servers"`
}

// Method is our struct used for storing methods and thier commands
type Method struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Servers []string `json:"servers"`
}

// Server is our struct used for storing a servers login information
type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
