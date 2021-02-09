package application

type AppServer struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Cors      bool   `json:"cors"`
	HttpLogs  bool   `json:"http_logs"`
	DbLogMode bool   `json:"db_log_mode"`
}

func NewAppServer() *AppServer {
	return &AppServer{}
}
