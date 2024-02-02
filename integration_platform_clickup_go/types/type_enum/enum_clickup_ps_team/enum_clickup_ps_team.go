package enum_clickup_ps_team

const (
	CONSULTING int = 0
	OFFSEC     int = 1
)

func ToString(value int) string {
	switch value {
	case CONSULTING:
		return "CONSULTING"
	case OFFSEC:
		return "OFFSEC"
	}
	return "unknown"
}
