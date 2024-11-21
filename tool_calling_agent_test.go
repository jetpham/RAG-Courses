package main

import "testing"

func Test_toolCallingAgent(t *testing.T) {
	setup, err := newSetup()
	if err != nil {
		t.Errorf("Error during setup: %v", err)
		return
	}
	tests := []struct {
		name    string
		prompt  string
		want    string
		wantErr bool
	}{
		{
			name:   "TestPhil",
			prompt: "What courses is Phil Peterson teaching in Fall 2024?",
			want: `
				Phil Peterson is teaching the following courses in Fall 2024:

				1. CS272-03 Software Development by Philip Peterson (In-Person, Tuesday Thursday, 2:40 PM - 4:25 PM, Building: LS, Room: G12).
				2. CS272-04 Software Development by Philip Peterson (In-Person, Tuesday Thursday, 8:00 AM - 9:45 AM, Building: LS, Room: G12).
				3. CS272L-01 Software Development Lab by Philip Peterson (In-Person, Wednesday, 1:00 PM - 2:30 PM, Building: MH, Room: 122).
				4. CS272L-02 Software Development Lab by Philip Peterson (In-Person, Wednesday, 2:55 PM - 4:25 PM, Building: MH, Room: 122).			
			`,
			wantErr: false,
		},
		{
			name:   "TestPHIL",
			prompt: "Which philosophy courses are offered this semester?",
			want: `
				The philosophy courses offered this semester are as follows:

				1. PHIL110-01 Great Philosophical Questions by Deena Lin (In-Person, MWF 09:15-10:20, KA 211, CRN: 41163).
				2. PHIL110-02 Great Philosophical Questions by Jea Oh (In-Person, MWF 01:00-02:05, LME 1225, CRN: 41164).
				3. PHIL110-03 Great Philosophical Questions by Jea Oh (In-Person, MWF 11:45-12:50, CO 107, CRN: 41165).
				4. PHIL110-04 Great Philosophical Questions by Richie Kim (Hybrid, MW 03:30-04:35, KA 311, CRN: 41166).
				5. PHIL110-05 Great Philosophical Questions by Purushottama Bilimoria (In-Person, TR 04:35-06:20, KA 263, CRN: 41167).
				6. PHIL110-06 Great Philosophical Questions by Purushottama Bilimoria (In-Person, TR 06:30-08:15, KA 263, CRN: 41168).
				7. PHIL195-02 FYS: Lovers of Wisdom by Thomas Cavanaugh (In-Person, TR 08:00-09:45, LM 365, CRN: 41170).
				8. PHIL202-01 Philosophy of Religion by Deena Lin (In-Person, MWF 01:00-02:05, LM 244A, CRN: 41172).
				9. PHIL204-01 Philosophy of Science by Krupa Patel (In-Person, MW 04:45-06:25, KA 163, CRN: 41174).
				10. PHIL204-02 Philosophy of Science by Krupa Patel (In-Person, MW 06:30-08:15, KA 163, CRN: 41175).
				11. PHIL205-01 Philosophy of Biology by Stephen Friesen (In-Person, MWF 10:30-11:35, ED 102, CRN: 41176).
				12. PHIL206-01 The Human Animal by Jennifer Fisher (In-Person, TR 02:40-04:25, ED 102, CRN: 42008).
				13. PHIL209-01 Aesthetics by Laurel Scotland-Stewart (In-Person, MWF 02:15-03:20, KA 111, CRN: 41177).
				14. PHIL209-02 Aesthetics by Laurel Scotland-Stewart (In-Person, MWF 03:30-04:35, KA 111, CRN: 41178).
				15. PHIL220-01 Asian Philosophy by Joshua Stoll (In-Person, MW 04:45-06:25, KA 263, CRN: 41179).
				16. PHIL240-02 Ethics by Joshua Carboni (In-Person, TR 02:40-04:25, LM 363, CRN: 41182).
				17. PHIL240-03 Ethics by Vida Pavesich (In-Person, MWF 09:15-10:20, KA 267, CRN: 41183).
				18. PHIL240-04 Ethics by Greig Mulberry (In-Person, MW 04:45-06:25, KA 267, CRN: 41184).
				19. PHIL240-05 Ethics by Greig Mulberry (In-Person, MW 06:30-08:15, KA 267, CRN: 41185).
				20. PHIL240-07 Ethics by Vida Pavesich (In-Person, MWF 01:00-02:05, KA 167, CRN: 41187).
				21. PHIL240-08 Ethics by Richie Kim (Hybrid, MW 02:15-03:20, KA 311, CRN: 41188).
				22. PHIL244-02 Environmental Ethics by Stephen Friesen (In-Person, MWF 09:15-10:20, ED 102, CRN: 41192).
				23. PHIL251-01 Mind, Freedom & Knowledge by Jennifer Fisher (In-Person, TR 04:35-06:20, ED 102, CRN: 41196).
				24. PHIL256-01 Existentialism by Brian Pines (In-Person, TR 09:55-11:40, ED 102, CRN: 41197).
				25. PHIL256-02 Existentialism by Brian Pines (In-Person, TR 12:45-02:30, ED 102, CRN: 41198).
				26. PHIL310-01 Ancient & Medieval Philosophy by Thomas Cavanaugh (In-Person, TR 09:55-11:40, LM 365, CRN: 41200).
				27. PHIL319-01 Logic by Nick Leonard (In-Person, MWF 11:45-12:50, ED 101, CRN: 41201).
				28. PHIL398-27 Contracts & First Nature by David Kim (In-Person, TBA, CRN: 42523).
				29. PHIL480-01 Topics in Contemporary Phil by (In-Person, MWF 10:30-11:35, ED 101, CRN: 41202).
				30. PHIL480-02 Topics in Contemporary Phil by David Kim (In-Person, TR 12:45-02:30, LM 363, CRN: 41203).
			`,
			wantErr: false,
		},
		{
			name:   "TestBio",
			prompt: "Where does Bioinformatics meet?",
			want: `
				Bioinformatics meets in the following courses:

				1. BIOL422-01 Bioinformatics by Naupaka Zimmerman (Meets MW 09:00 - 10:15, Building: KA, Room: 311).
				2. BTEC640-01 Bioinformatics by Naupaka Zimmerman (Meets MW 16:45 - 18:25, Building: KA, Room: 111).
				3. BTEC640-02 Bioinformatics (Meets MW 16:45 - 18:25, Building: HR, Room: 136).
			`,
			wantErr: false,
		},
		{
			name:   "TestGuitar",
			prompt: "Can I learn guitar this semester?",
			want: `
				Yes, you can learn guitar this semester.

				cited courses:
				1. MUS121-01 Guitar and Bass Lessons by Christopher Ruscoe (In-Person, 8/20/24 - 11/28/24)
				2. MUS121-02 Guitar and Bass Lessons by Christopher Ruscoe (In-Person, 8/20/24 - 11/28/24)
			`,
			wantErr: false,
		},
		{
			name:   "TestMultiple",
			prompt: "I would like to take a Rhetoric course from Phil Choong. What can I take?",
			want: `
				You can take the following Rhetoric courses from Phil Choong:

				1. RHET103-05 Public Speaking by Philip Choong (In-Person, MWF, 10:30 AM - 11:35 AM, Building LM, Room 346A)
				2. RHET103-26 Public Speaking by Philip Choong (In-Person, MWF, 11:45 AM - 12:50 PM, Building LM, Room 346A)
				3. RHET195-02 FYS: Podcasts: Eloquentia & Aud by Philip Choong (In-Person, MWF, 2:15 PM - 3:20 PM, Building LM, Room 352)
				4. RHET328-01 Speaking Center Internship by Philip Choong (In-Person, Tuesday, 4:35 PM - 6:25 PM, Building LM, Room 345)
			`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toolCallingAgent(setup, tt.prompt)
			similar, reason := isSimilar(setup, got, tt.want)
			if !similar {
				t.Errorf("toolCallingAgent() = \n\n%v\n\nwant:\n\n%v\n\nreason:\n\n%v", got, tt.want, reason)
			}
		})
	}
}
