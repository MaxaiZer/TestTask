package unit

type MockNotifyService struct {
	notifiedEmails []string
}

func (s *MockNotifyService) IsNotified(email string) bool {

	for _, _email := range s.notifiedEmails {
		if _email == email {
			return true
		}
	}

	return false
}

func (s *MockNotifyService) NotifyAboutIpChange(email string) {
	s.notifiedEmails = append(s.notifiedEmails, email)
}
