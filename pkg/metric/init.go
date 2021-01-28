package metric

func init() {
	setGlobal(&noOpMetric{})
}
