package v1alpha1

const (
	IgnoreUnderlyingDeletionAnnotation = "github.sanjivmadhavan.io/ignore-underlying-deletion"
)

func (s *Space) HasIgnoreUnderlyingDeletionAnnotation() bool {
	if _, ok := s.Annotations[IgnoreUnderlyingDeletionAnnotation]; ok {
		return true
	}

	return false
}
