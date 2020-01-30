
// 1.contact
func getcontact(mail string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from contact where email='%s'", mail))
	if err != nil {
		panic(err.Error())
	}
	contactInfo := []string{"", "", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2], &userInfo[3])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return contactInfo, err
}

//2.Client
func getclient(contact_id string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from client where email='%s'", contact_id))
	if err != nil {
		panic(err.Error())
	}
	clientInfo := []string{"", "", "", "", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2], &userInfo[3], &userInfo[4], &userInfo[5])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return clientInfo, err
}

//3.service
func getservice(string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from service where email='%s'"))
	if err != nil {
		panic(err.Error())
	}
	serviceInfo := []string{"", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return serviceInfo, err
}

//4.user
func getuser(contact_id string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from user where email='%s'", contact_id))
	if err != nil {
		panic(err.Error())
	}
	userInfo := []string{"", "", "", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2], &userInfo[3], &userInfo[4])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return userInfo, err
}

//5.client_service
func getclient_service(client_id string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from client_service where email='%s'", client_id))
	if err != nil {
		panic(err.Error())
	}
	client_serviceInfo := []string{"", "", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2], &userInfo[3])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return client_serviceInfo, err
}

//6.medical
func getmedical(user_id string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from medical where email='%s'", user_id))
	if err != nil {
		panic(err.Error())
	}
	medicalInfo := []string{"", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return medicalInfo, err
}

//7.activity
func getactivity(string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from activity where email='%s'"))
	if err != nil {
		panic(err.Error())
	}
	activityInfo := []string{"", "", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2], &userInfo[3])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return activityInfo, err
}

//8.`user_activity
func getuser_activity(user_id string) ([]string, error) {
	rows, err := c.DB.Query(fmt.Sprintf("select * from user_activity where email='%s'", user_id))
	if err != nil {
		panic(err.Error())
	}
	user_activityInfo := []string{"", "", "", ""}
	for rows.Next() {
		err = rows.Scan(&userInfo[0], &userInfo[1], &userInfo[2], &userInfo[3])
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	return user_activityInfo, err
}

//9.client_form
