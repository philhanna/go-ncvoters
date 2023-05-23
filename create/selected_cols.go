package create

var selectedCols = []string{
	"county_id",
	"voter_reg_num",
	"last_name",
	"first_name",
	"middle_name",
	"name_suffix_lbl",
	"status_cd",
	"reason_cd",
	"res_street_address",
	"res_city_desc",
	"state_cd",
	"zip_code",
	"full_phone_number",
	"race_code",
	"ethnic_code",
	"party_cd",
	"gender_code",
	"birth_year",
	"age_at_year_end",
	"birth_state",
}

var sanitizeCols = []string{
	"last_name",
	"first_name",
	"middle_name",
	"res_street_address",
	"res_city_desc",
}
