package service

// service层方法，接收请求结构体或特定参数执行dao方法

func (svc *RuleService) GetAllEmails() ([]string, error) {
	all, err := svc.dao.ListEmails()
	if err != nil {
		return nil, err
	}
	var emails []string
	for _, email := range all {
		emails = append(emails, email.Address)
	}
	return emails, nil
}
