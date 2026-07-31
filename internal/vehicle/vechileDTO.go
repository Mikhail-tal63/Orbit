package vehicle

type CreatVechileRequest struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        int32  `json:"year"`
	Color       string `json:"color"`
	PlateNumber string `json:"plate_number"`
}
