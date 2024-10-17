package dependent

import "regexp"

type CheckPartnerData struct{}

func (e *CheckPartnerData) IsEmpty(data string) bool {
	return data == ""
}

func (e *CheckPartnerData) IsEmailValidFormat(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (e *CheckPartnerData) IsPhoneValidFormat(phone string) bool {
	re := regexp.MustCompile(`^09\d{8}$`)
	return re.MatchString(phone)
}

func (e *CheckPartnerData) IsSexValidFormat(sex string) bool {
	return sex == "" || sex == "男" || sex == "女"
}

func (e *CheckPartnerData) IsSocialMediaOptionOne(instagram, facebook, threads, twitter string) bool {
	return instagram != "" || facebook != "" || threads != "" || twitter != ""
}