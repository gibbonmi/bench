package specbuild

// Status returns the durable compact projection for slug.
func (s *Service) Status(slug string) (Status, error) {
	if _, err := s.resolve(slug); err != nil {
		return Status{}, err
	}
	run, found, err := s.load(slug)
	if err != nil {
		return Status{}, err
	}
	if !found {
		return Status{Slug: slug, State: "empty", Next: "bench spec build start " + slug}, nil
	}
	return run.status(), nil
}
