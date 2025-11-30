package v2

// هذا الملف مؤقت للإصدار المستقبلي v2
// سيتم تطويره لاحقاً عندما نحتاج لإصدار جديد من API

// APIV2Info معلومات عن الإصدار المستقبلي
type APIV2Info struct {
	Version     string `json:"version"`
	Status      string `json:"status"`
	Description string `json:"description"`
	PlannedRelease string `json:"planned_release"`
}

// GetV2Info معلومات عن الإصدار v2
func GetV2Info() APIV2Info {
	return APIV2Info{
		Version:       "v2.0.0",
		Status:        "planned",
		Description:   "الإصدار الثاني من NawthTech API مع ميزات متقدمة",
		PlannedRelease: "Q4 2024",
	}
}