package book

import "context"

func (s *Service) LoadCache(ctx context.Context, limit int) error {
	books, err := s.repository.GetAllWithLimit(ctx, limit)
	if err != nil {
		s.logger.Error("cache load error", "err", err)
		return err
	}

	for _, b := range books {
		if err := s.cache.Set(ctx, b.ID.String(), b); err != nil {
			s.logger.Error("cache set error", "err", err)
		}
	}
	s.logger.Info("cache load", "len", len(books))
	return nil
}
