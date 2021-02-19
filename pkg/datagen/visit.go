package gen

type Visit struct {
	LocationName string `json:"loc_name"`
	LocationCode string `json:"loc_code"`
	Service      string `json:"service"`
	PatientClass string `json:"pt_class"`
	PatientType  string `json:"pt_type"`
	AdmitType    string `json:"admit_type"`
	AdmitSource  string `json:"admit_source"`
}
