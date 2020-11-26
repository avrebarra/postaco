package config

type ServerCode string

type ServerResponse struct {
	Status  ServerCode  `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ServerResponseBasic struct {
	ErrorTitle   string `json:"error_title"`
	ErrorContent string `json:"error_content"`
}

type ServerResponseName int

const (
	ServerRespSuccess ServerResponseName = iota
	ServerRespErrUnauthorizedInternal
	ServerRespErrPathNotFound
	ServerRespErrUnexpected
	ServerRespErrValidation
	ServerRespErrProcessFailed
)

var HashServerResponse map[ServerResponseName]ServerResponse = map[ServerResponseName]ServerResponse{
	ServerRespSuccess: {
		Status:  "00",
		Message: "success",
	},

	ServerRespErrPathNotFound: {
		Status:  "03",
		Message: "path not found",
		Data: ServerResponseBasic{
			ErrorTitle:   "Terjadi Kesalahan",
			ErrorContent: "Resource yang dimaksud tidak ditemukan.",
		},
	},

	ServerRespErrUnexpected: {
		Status:  "05",
		Message: "unexpected error",
		Data: ServerResponseBasic{
			ErrorTitle:   "Terjadi Kesalahan",
			ErrorContent: "Permintaan gagal. Silakan hubungi Agent Helpdesk.",
		},
	},

	ServerRespErrValidation: {
		Status:  "06",
		Message: "validation error",
		Data: ServerResponseBasic{
			ErrorTitle:   "Validasi Gagal",
			ErrorContent: "Proses validasi input gagal. Pastikan data yang anda masukkan sudah benar.",
		},
	},

	ServerRespErrProcessFailed: {
		Status:  "07",
		Message: "process failed",
		Data: ServerResponseBasic{
			ErrorTitle:   "Proses Gagal",
			ErrorContent: "Gagal menyelesaikan proses. Silakan hubungi Agent Helpdesk.",
		},
	},
}
