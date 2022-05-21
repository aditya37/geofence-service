package tile38

type (
	Field struct {
		Timestamp float64
		Speed     float64
	}
	SetKey struct {
		Key      string
		ObjectId string
		Fields   Field
		Lat      float64
		Long     float64
	}
)
