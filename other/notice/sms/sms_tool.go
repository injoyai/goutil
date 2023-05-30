package sms

// dealPhone 处理电话
func dealPhone(phone []string) []string {
	for i, v := range phone {
		if len(v) == 11 {
			phone[i] = "+86" + v
		}
	}
	return phone
}
