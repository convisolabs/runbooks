package enum_clickup_ps_team

const (
	CONSULTING           int = 0
	OFFSEC               int = 1
	EDUCATIONANDCOMMUNIY int = 2
)

func ToString(value int) string {
	switch value {
	case CONSULTING:
		return "CONSULTING"
	case OFFSEC:
		return "OFFSEC"
	case EDUCATIONANDCOMMUNIY:
		return "Education and Community"
	}
	return "unknown"
}
