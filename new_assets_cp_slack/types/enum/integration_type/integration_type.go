package enum_integration_type

const (
	CRAWLER_FORTIFY int = 0
	ASSETS_CP       int = 1
)

func ToString(value int) string {
	switch value {
	case CRAWLER_FORTIFY:
		return "CRAWLER_FORTIFY"
	case ASSETS_CP:
		return "ASSETS_CP"
	}
	return "unknown"

}
