


	err := row.Scan(&id,&l.Name,&l.email,&services,&description,&l.created_at,&l.modified_at,&address,
 &l.primary_ph,&l.secondary_ph)
	if err != nil {
		return nil, fmt.Errorf("client not found: %v", err)
	}
	err = json.Unmarshal([]byte(Address), &l.Address)
	if err != nil {
		fmt.Printf("unmarshalling address failed: %v", err)
		fmt.Println(l.Address)
	}

  err = json.Unmarshal([]byte( services), &l.services)
	if err != nil {
		fmt.Printf("unmarshalling  services failed: %v", err)
		fmt.Println(l.services)
	}

  err = json.Unmarshal([]byte(description), &l.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(l.description)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	l.Stats = *stats
	return &l, err
}

//3.service
func (ms *dbStorage) getservice(email string) (*atium.serviceInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
	"select S.id,C.client.id,S.name,S.description,S.price,",
  "C.enabledAt,C.enabled from service S",
  "LEFT JOIN client_service C ON C.service_id = S.id",
  "WHERE email = ?")

	var description string
	var id int64
	u := atium.serviceInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&u.client_id,&u.name,&description,&u.price,
  &u.enabledAt,&u.enabled )
	if err != nil {
		return nil, fmt.Errorf("service not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &u.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(u.description)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//4.client_service
func (
		return nil, fmt.Errorf("client_service not found: %v", err)
	}
	err = json.Unmarshal([]byte(description), &u.description)
	if err != nil {
		fmt.Printf("unmarshalling description failed: %v", err)
		fmt.Println(u.description)
	}
  err = json.Unmarshal([]byte(services), &u.services)
	if err != nil {
		fmt.Printf("unmarshalling services failed: %v", err)
		fmt.Println(u.services)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}

//5.medical

func (ms *dbStorage) getmedical(email string) (*atium.medicalInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
"select M.id,U.name,U.gender,U.contact_id,U.Dob,",
"U.created_at,M.medical_record from medical M",
"LEFT JOIN user U ON U.id= M.user_id",
"WHERE email = ?")
	var medical_record string
	var id int64
	x := atium.medicalInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&id,&x.name,&x.gender,&x.contact_id,&x.Dob,
  &x.created_at,&medical_record)
	if err != nil {
		return nil, fmt.Errorf("medical_record not found: %v", err)
	}
	err = json.Unmarshal([]byte(medical_record), &x.medical_record)
	if err != nil {
		fmt.Printf("unmarshalling medical_record failed: %v", err)
		fmt.Println(x.medical_record)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	x.Stats = *stats
	return &x, err
}

//6.survey
func (ms *dbStorage) getsurvey(email string) (*atium.surveyInfo, error) {
	qs := fmt.Sprintf("%s %s %s %s",
	"select S.id,U.name,U.gender,U.dob,U.created_at,",
  "S.client_form_map_id,S.form_entry_id from survey S"
  "LEFT JOIN user U ON U.id= S.user_id",
  "WHERE email = ?")

	var id int64
	u := atium.surveyInfo{}
	row := ms.db.QueryRow(qs, email)

	err := row.Scan(&u.id,&u.name,&u.gender,&u.dob,&u.created_at,
  &u.client_form_map_id,&u.form_entry_id)
	if err != nil {
		return nil, fmt.Errorf("survey not found: %v", err)
	}
	stats, err := getLatestStats(ms.db, id)
	if err != nil {
		return nil, fmt.Errorf("err getting stats: %v", err)
	}
	u.Stats = *stats
	return &u, err
}
