package pipelineprocessor

// k8sObjectName returns a valid Kubernetes object name by combining two strings
func k8sObjectName(str1, str2 string) string {
	combined := str1 + "-" + str2
	if len(combined) > 60 {
		return combined[:60]
	}
	return combined
}
