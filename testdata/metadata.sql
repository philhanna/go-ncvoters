PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE columns (
  name           TEXT,
  dataType       TEXT,
  description    TEXT
);
INSERT INTO columns VALUES('county_id','int','County identification number');
INSERT INTO columns VALUES('county_desc','varchar(15)','County name');
INSERT INTO columns VALUES('voter_reg_num','char(12)','Voter registration number (unique to county)');
INSERT INTO columns VALUES('ncid','char(12)','North Carolina identification (NCID) number');
INSERT INTO columns VALUES('last_name','varchar(25)','Voter last name');
INSERT INTO columns VALUES('first_name','varchar(20)','Voter first name');
INSERT INTO columns VALUES('middle_name','varchar(20)','Voter middle name');
INSERT INTO columns VALUES('name_suffix_lbl','char(3)','Voter suffix name (JR, III, etc.)');
INSERT INTO columns VALUES('status_cd','char(2)','Registration status code');
INSERT INTO columns VALUES('voter_status_desc','varchar(25)','Registration status description');
INSERT INTO columns VALUES('reason_cd','varchar(2)','Registration status reason code');
INSERT INTO columns VALUES('voter_status_reason_desc','varchar(60)','Registration status reason description');
INSERT INTO columns VALUES('res_street_address','varvarchar(65)','Residential street address');
INSERT INTO columns VALUES('res_city_desc','varchar(60)','Residential city name');
INSERT INTO columns VALUES('state_cd','varchar(2)','Residential address state code');
INSERT INTO columns VALUES('zip_code','char(9)','Residential address zip code');
INSERT INTO columns VALUES('mail_addr1','varchar(40)','Mailing address line 1');
INSERT INTO columns VALUES('mail_addr2','varchar(40)','Mailing address line 2');
INSERT INTO columns VALUES('mail_addr3','varchar(40)','Mailing address line 3');
INSERT INTO columns VALUES('mail_addr4','varchar(40)','Mailing address line 4');
INSERT INTO columns VALUES('mail_city','varchar(3)0','Mailing address city name');
INSERT INTO columns VALUES('mail_state','varchar(2)','Mailing address city code');
INSERT INTO columns VALUES('mail_zipcode','char(9)','Mailing address zip code');
INSERT INTO columns VALUES('full_phone_number','varchar(12)','Full phone number including area code');
INSERT INTO columns VALUES('confidential_ind','char(1)','Confidential indicator (by General Statute certain data is confidential for this record)');
INSERT INTO columns VALUES('registr_dt','char(10)','Registration date');
INSERT INTO columns VALUES('race_code','char(3)','Race code');
INSERT INTO columns VALUES('ethnic_code','char(3)','Ethnicity code');
INSERT INTO columns VALUES('party_cd','char(3)','Registered party code');
INSERT INTO columns VALUES('gender_code','char(1)','Gender/sex code');
INSERT INTO columns VALUES('birth_year','char(4)','Year of birth');
INSERT INTO columns VALUES('age_at_year_end','char(3)','Age at end of the year (was: birth_age - 02/09/2022)');
INSERT INTO columns VALUES('birth_state','varchar(2)','Birth state');
INSERT INTO columns VALUES('drivers_lic','char(1)','Drivers license (Y/N)');
INSERT INTO columns VALUES('precinct_abbrv','varchar(6)','Precinct abbreviation');
INSERT INTO columns VALUES('precinct_desc','varchar(60)','Precinct name');
INSERT INTO columns VALUES('municipality_abbrv','varchar(6)','Municipality jurisdiction abbreviation');
INSERT INTO columns VALUES('municipality_desc','varchar(60)','Municipality jurisdiction name');
INSERT INTO columns VALUES('ward_abbrv','varchar(6)','Ward jurisdiction abbreviation');
INSERT INTO columns VALUES('ward_desc','varchar(60)','Ward jurisdiction name');
INSERT INTO columns VALUES('cong_dist_abbrv','varchar(6)','Congressional district abbreviation');
INSERT INTO columns VALUES('super_court_abbrv','varchar(6)','Superior court jurisdiction abbreviation');
INSERT INTO columns VALUES('judic_dist_abbrv','varchar(6)','Judicial district abbreviation');
INSERT INTO columns VALUES('nc_senate_abbrv','varchar(6)','NC Senate jurisdiction abbreviation');
INSERT INTO columns VALUES('nc_house_abbrv','varchar(6)','NC House jurisdiction abbreviation');
INSERT INTO columns VALUES('county_commiss_abbrv','varchar(6)','County commisioner jurisdiction abbreviation');
INSERT INTO columns VALUES('county_commiss_desc','varchar(60)','County commisioner jurisdiction name');
INSERT INTO columns VALUES('township_abbrv','varchar(6)','Township jurisdiction abbreviation');
INSERT INTO columns VALUES('township_desc','varchar(60)','Township jurisdiction name');
INSERT INTO columns VALUES('school_dist_abbrv','varchar(6)','School district abbreviation');
INSERT INTO columns VALUES('school_dist_desc','varchar(60)','School district name');
INSERT INTO columns VALUES('fire_dist_abbrv','varchar(6)','Fire district abbreviation');
INSERT INTO columns VALUES('fire_dist_desc','varchar(60)','Fir district name');
INSERT INTO columns VALUES('water_dist_abbrv','varchar(6)','Water district abbreviation');
INSERT INTO columns VALUES('water_dist_desc','varchar(60)','Water district name');
INSERT INTO columns VALUES('sewer_dist_abbrv','varchar(6)','Sewer district abbreviation');
INSERT INTO columns VALUES('sewer_dist_desc','varchar(60)','Sewer district name');
INSERT INTO columns VALUES('sanit_dist_abbrv','varchar(6)','Sanitation district abbreviation');
INSERT INTO columns VALUES('sanit_dist_desc','varchar(60)','Sanitation district name');
INSERT INTO columns VALUES('rescue_dist_abbrv','varchar(6)','Rescue district abbreviation');
INSERT INTO columns VALUES('rescue_dist_desc','varchar(60)','Rescue district name');
INSERT INTO columns VALUES('munic_dist_abbrv','varchar(6)','Municpal district abbreviation');
INSERT INTO columns VALUES('munic_dist_desc','varchar(60)','Municipal district name');
INSERT INTO columns VALUES('dist_1_abbrv','varchar(6)','Presecutorial district abbreviation');
INSERT INTO columns VALUES('dist_1_desc','varchar(60)','Presecutorial district name');
INSERT INTO columns VALUES('vtd_abbrv','varchar(6)','Voter tabulation district abbreviation');
INSERT INTO columns VALUES('vtd_desc','varchar(60)','Voter tabulation district name');
CREATE TABLE status_codes (
  status         TEXT,
  description    TEXT
);
INSERT INTO status_codes VALUES('A','ACTIVE');
INSERT INTO status_codes VALUES('D','DENIED');
INSERT INTO status_codes VALUES('I','INACTIVE');
INSERT INTO status_codes VALUES('R','REMOVED');
INSERT INTO status_codes VALUES('S','TEMPORARY (APPLICABLE TO MILITARY AND OVERSEAS)');
CREATE TABLE race_codes (
  race           TEXT,
  description    TEXT
);
INSERT INTO race_codes VALUES('A','ASIAN');
INSERT INTO race_codes VALUES('B','BLACK or AFRICAN AMERICAN');
INSERT INTO race_codes VALUES('I','AMERICAN INDIAN or ALASKA NATIVE');
INSERT INTO race_codes VALUES('M','TWO or MORE RACES ');
INSERT INTO race_codes VALUES('O','OTHER');
INSERT INTO race_codes VALUES('P','NATIVE HAWAIIAN or PACIFIC ISLANDER');
INSERT INTO race_codes VALUES('U','UNDESIGNATED');
INSERT INTO race_codes VALUES('W','WHITE');
CREATE TABLE ethnic_codes (
  ethnicity      TEXT,
  description    TEXT
);
INSERT INTO ethnic_codes VALUES('HL','HISPANIC or LATINO');
INSERT INTO ethnic_codes VALUES('NL','NOT HISPANIC or NOT LATINO');
INSERT INTO ethnic_codes VALUES('UN','UNDESIGNATED');
CREATE TABLE county_codes (
  county_id      TEXT,
  county         TEXT
);
INSERT INTO county_codes VALUES('0','YANCEY');
INSERT INTO county_codes VALUES('1','ALAMANCE');
INSERT INTO county_codes VALUES('2','ALEXANDER');
INSERT INTO county_codes VALUES('3','ALLEGHANY');
INSERT INTO county_codes VALUES('4','ANSON');
INSERT INTO county_codes VALUES('5','ASHE');
INSERT INTO county_codes VALUES('6','AVERY');
INSERT INTO county_codes VALUES('7','BEAUFORT');
INSERT INTO county_codes VALUES('8','BERTIE');
INSERT INTO county_codes VALUES('9','BLADEN');
INSERT INTO county_codes VALUES('10','BRUNSWICK');
INSERT INTO county_codes VALUES('11','BUNCOMBE');
INSERT INTO county_codes VALUES('12','BURKE');
INSERT INTO county_codes VALUES('13','CABARRUS');
INSERT INTO county_codes VALUES('14','CALDWELL');
INSERT INTO county_codes VALUES('15','CAMDEN');
INSERT INTO county_codes VALUES('16','CARTERET');
INSERT INTO county_codes VALUES('17','CASWELL');
INSERT INTO county_codes VALUES('18','CATAWBA');
INSERT INTO county_codes VALUES('19','CHATHAM');
INSERT INTO county_codes VALUES('20','CHEROKEE');
INSERT INTO county_codes VALUES('21','CHOWAN');
INSERT INTO county_codes VALUES('22','CLAY');
INSERT INTO county_codes VALUES('23','CLEVELAND');
INSERT INTO county_codes VALUES('24','COLUMBUS');
INSERT INTO county_codes VALUES('25','CRAVEN');
INSERT INTO county_codes VALUES('26','CUMBERLAND');
INSERT INTO county_codes VALUES('27','CURRITUCK');
INSERT INTO county_codes VALUES('28','DARE');
INSERT INTO county_codes VALUES('29','DAVIDSON');
INSERT INTO county_codes VALUES('30','DAVIE');
INSERT INTO county_codes VALUES('31','DUPLIN');
INSERT INTO county_codes VALUES('32','DURHAM');
INSERT INTO county_codes VALUES('33','EDGECOMBE');
INSERT INTO county_codes VALUES('34','FORSYTH');
INSERT INTO county_codes VALUES('35','FRANKLIN');
INSERT INTO county_codes VALUES('36','GASTON');
INSERT INTO county_codes VALUES('37','GATES');
INSERT INTO county_codes VALUES('38','GRAHAM');
INSERT INTO county_codes VALUES('39','GRANVILLE');
INSERT INTO county_codes VALUES('40','GREENE');
INSERT INTO county_codes VALUES('41','GUILFORD');
INSERT INTO county_codes VALUES('42','HALIFAX');
INSERT INTO county_codes VALUES('43','HARNETT');
INSERT INTO county_codes VALUES('44','HAYWOOD');
INSERT INTO county_codes VALUES('45','HENDERSON');
INSERT INTO county_codes VALUES('46','HERTFORD');
INSERT INTO county_codes VALUES('47','HOKE');
INSERT INTO county_codes VALUES('48','HYDE');
INSERT INTO county_codes VALUES('49','IREDELL');
INSERT INTO county_codes VALUES('50','JACKSON');
INSERT INTO county_codes VALUES('51','JOHNSTON');
INSERT INTO county_codes VALUES('52','JONES');
INSERT INTO county_codes VALUES('53','LEE');
INSERT INTO county_codes VALUES('54','LENOIR');
INSERT INTO county_codes VALUES('55','LINCOLN');
INSERT INTO county_codes VALUES('56','MACON');
INSERT INTO county_codes VALUES('57','MADISON');
INSERT INTO county_codes VALUES('58','MARTIN');
INSERT INTO county_codes VALUES('59','MCDOWELL');
INSERT INTO county_codes VALUES('60','MECKLENBURG');
INSERT INTO county_codes VALUES('61','MITCHELL');
INSERT INTO county_codes VALUES('62','MONTGOMERY');
INSERT INTO county_codes VALUES('63','MOORE');
INSERT INTO county_codes VALUES('64','NASH');
INSERT INTO county_codes VALUES('65','NEWHANOVER');
INSERT INTO county_codes VALUES('66','NORTHAMPTON');
INSERT INTO county_codes VALUES('67','ONSLOW');
INSERT INTO county_codes VALUES('68','ORANGE');
INSERT INTO county_codes VALUES('69','PAMLICO');
INSERT INTO county_codes VALUES('70','PASQUOTANK');
INSERT INTO county_codes VALUES('71','PENDER');
INSERT INTO county_codes VALUES('72','PERQUIMANS');
INSERT INTO county_codes VALUES('73','PERSON');
INSERT INTO county_codes VALUES('74','PITT');
INSERT INTO county_codes VALUES('75','POLK');
INSERT INTO county_codes VALUES('76','RANDOLPH');
INSERT INTO county_codes VALUES('77','RICHMOND');
INSERT INTO county_codes VALUES('78','ROBESON');
INSERT INTO county_codes VALUES('79','ROCKINGHAM');
INSERT INTO county_codes VALUES('80','ROWAN');
INSERT INTO county_codes VALUES('81','RUTHERFORD');
INSERT INTO county_codes VALUES('82','SAMPSON');
INSERT INTO county_codes VALUES('83','SCOTLAND');
INSERT INTO county_codes VALUES('84','STANLY');
INSERT INTO county_codes VALUES('85','STOKES');
INSERT INTO county_codes VALUES('86','SURRY');
INSERT INTO county_codes VALUES('87','SWAIN');
INSERT INTO county_codes VALUES('88','TRANSYLVANIA');
INSERT INTO county_codes VALUES('89','TYRRELL');
INSERT INTO county_codes VALUES('90','UNION');
INSERT INTO county_codes VALUES('91','VANCE');
INSERT INTO county_codes VALUES('92','WAKE');
INSERT INTO county_codes VALUES('93','WARREN');
INSERT INTO county_codes VALUES('94','WASHINGTON');
INSERT INTO county_codes VALUES('95','WATAUGA');
INSERT INTO county_codes VALUES('96','WAYNE');
INSERT INTO county_codes VALUES('97','WILKES');
INSERT INTO county_codes VALUES('98','WILSON');
INSERT INTO county_codes VALUES('99','YADKIN');
COMMIT;