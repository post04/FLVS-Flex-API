package schoolapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/saihon/saihon"
)

var (
	transport = http.Transport{}
)

func login(username, password string) (string, string, error) {
	req, err := http.NewRequest("POST", "https://login.flvs.net/", strings.NewReader(fmt.Sprintf("ReturnUrl=&Username=%v&Password=%v&RememberUsername=false", username, password)))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", err
	}
	//fmt.Println("(Login) -", resp.Status)
	return strings.Split(strings.Split(resp.Header.Get("Set-Cookie"), "=")[1], ";")[0], resp.Header.Get("Location"), nil
}

func secondStep(membershipURL string) (string, string, string, error) {
	var aspx, asp string
	req, err := http.NewRequest("GET", membershipURL, nil)
	if err != nil {
		return "", "", "", err
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", "", err
	}
	//fmt.Println("(Second) -", resp.Status)
	for _, cookie := range resp.Header["Set-Cookie"] {
		if strings.HasPrefix(cookie, "ASP") {
			asp = strings.Split(strings.Split(cookie, "=")[1], ";")[0]
		} else if strings.HasPrefix(cookie, ".ASPX") {
			aspx = strings.Split(strings.Split(cookie, "=")[1], ";")[0]
		}
	}
	return aspx, asp, resp.Header.Get("Location"), nil
}

func thirdStep(defaultURL, cookieToSubmit string) (string, string, string, error) {
	var postData, rpts string
	req, err := http.NewRequest("GET", defaultURL, nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Add("Cookie", cookieToSubmit)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", "", err
	}
	//fmt.Println("(Third) -", resp.Status)
	document, err := saihon.Parse(resp.Body)
	if err != nil {
		//fmt.Println(err)
		return "", "", "", err
	}
	link := document.Body().Node.FirstChild.Attr[2].Val
	wresult := document.Body().Node.FirstChild.FirstChild.NextSibling.Attr[2].Val
	wctx := document.Body().Node.FirstChild.FirstChild.NextSibling.NextSibling.Attr[2].Val
	postData = fmt.Sprintf("wa=%v&wresult=%v&wctx=%v", "wsignin1.0", url.QueryEscape(wresult), url.QueryEscape(wctx))
	rpts = strings.Split(strings.Join(strings.Split(resp.Header.Get("Set-Cookie"), "=")[1:], "="), ";")[0]
	return postData, rpts, link, nil
}

func fourthStep(avalibleSitesURL, cookieToSubmit, postData string) (string, error) {
	req, err := http.NewRequest("POST", avalibleSitesURL, strings.NewReader(postData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", cookieToSubmit)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	//fmt.Println("(Fourth) -", resp.Status)
	return strings.Split(strings.Split(resp.Header.Get("Set-Cookie"), "=")[1], ";")[0], nil
}

func fifthStep() (string, error) {
	req, err := http.NewRequest("GET", "https://vsa.flvs.net/", nil)
	if err != nil {
		return "", err
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	//fmt.Println("(Fifth) -", resp.Status)
	return resp.Header.Get("Location"), nil
}

func sixthStep(URL, cookieToGive string) (string, string, string, error) {
	var postData, rpts string
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Add("Cookie", cookieToGive)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", "", err
	}
	//fmt.Println("(Sixth) -", resp.Status)
	document, err := saihon.Parse(resp.Body)
	if err != nil {
		//fmt.Println(err)
		return "", "", "", err
	}
	link := document.Body().Node.FirstChild.Attr[2].Val
	wresult := document.Body().Node.FirstChild.FirstChild.NextSibling.Attr[2].Val
	wctx := document.Body().Node.FirstChild.FirstChild.NextSibling.NextSibling.Attr[2].Val
	postData = fmt.Sprintf("wa=%v&wresult=%v&wctx=%v", "wsignin1.0", url.QueryEscape(wresult), url.QueryEscape(wctx))
	rpts = strings.Split(strings.Join(strings.Split(resp.Header.Get("Set-Cookie"), "=")[1:], "="), ";")[0]
	return postData, rpts, link, nil
}

func seventhStep(postData, URL string) (string, string, error) {
	req, err := http.NewRequest("POST", URL, strings.NewReader(postData))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", err
	}
	//fmt.Println("(Seventh) -", resp.Status)
	return strings.Split(strings.Split(resp.Header.Get("Set-Cookie"), "=")[1], ";")[0], "https://vsa.flvs.net/", nil
}

func eigthStep(URL, cookieToGive string) (string, string, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Cookie", cookieToGive)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", err
	}
	//fmt.Println("(Eigth) -", resp.Status)
	return strings.Split(strings.Split(resp.Header.Get("Set-Cookie"), "=")[1], ";")[0], "https://sts.flvs.net" + resp.Header.Get("Location"), nil
}

func ninethStep(URL, cookieToGive string) (string, string, string, error) {
	var svSID, setSEC string
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Add("Cookie", cookieToGive)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", "", "", err
	}
	//fmt.Println("(Nineth) -", resp.Status)
	for _, cookie := range resp.Header["Set-Cookie"] {
		if strings.HasPrefix(cookie, "SVSID") {
			svSID = strings.Split(strings.Split(cookie, "=")[1], ";")[0]
		} else if strings.HasPrefix(cookie, "SetSec") {
			setSEC = strings.Split(cookie, ";")[0] + ";"
		}
	}
	return svSID, setSEC, "https://vsa.flvs.net/Resources/JavaScript/api.aspx?noext=1&noext=1", nil
}

func tenthStep(URL, cookieToGive string) (string, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", cookieToGive)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	//fmt.Println("(Tenth) -", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseAPIKey(html string) string {
	toReturn := strings.Split(html, "\": \"")[1]
	toReturn = strings.Split(toReturn, "\n")[0]
	toReturn = strings.ReplaceAll(toReturn, "\" }", "")
	toReturn = strings.ReplaceAll(toReturn, "\r", "")
	return toReturn
}

// CourseInfo is the entire response from the FLVS API that contains loads of information
type CourseInfo struct {
	Data []struct {
		ReportingCourseSegmentID    int         `json:"ReportingCourseSegmentId"`
		CourseVersionID             int         `json:"CourseVersionId"`
		CourseID                    int         `json:"CourseId"`
		CourseName                  string      `json:"CourseName"`
		ExternalCourseName          string      `json:"ExternalCourseName"`
		CourseSegment               int         `json:"CourseSegment"`
		CourseCode                  string      `json:"CourseCode"`
		CourseType                  string      `json:"CourseType"`
		RequestDate                 string      `json:"RequestDate"`
		PreferredStartDateString    string      `json:"PreferredStartDateString"`
		TeacherID                   int         `json:"TeacherId"`
		TeacherName                 string      `json:"TeacherName"`
		TeacherEmailAddress         string      `json:"TeacherEmailAddress"`
		TeacherPhoneNumber          string      `json:"TeacherPhoneNumber"`
		PercentComplete             float64     `json:"PercentComplete"`
		ActiveDate                  string      `json:"ActiveDate"`
		AverageGradeString          string      `json:"AverageGradeString"`
		FinalGradeString            string      `json:"FinalGradeString"`
		AverageLetterGrade          string      `json:"AverageLetterGrade"`
		ShouldShowLetterGrade       bool        `json:"ShouldShowLetterGrade"`
		NumericGrade                interface{} `json:"NumericGrade"`
		LetterGrade                 interface{} `json:"LetterGrade"`
		VirtualSchoolName           string      `json:"VirtualSchoolName"`
		LmsImportDate               string      `json:"LmsImportDate"`
		PhysicalSchoolID            int         `json:"PhysicalSchoolId"`
		CourseWorkLocation          string      `json:"CourseWorkLocation"`
		ClientName                  interface{} `json:"ClientName"`
		InvoiceNumber               interface{} `json:"InvoiceNumber"`
		InvoicedByName              interface{} `json:"InvoicedByName"`
		TuitionFee                  float64     `json:"TuitionFee"`
		Balance                     float64     `json:"Balance"`
		FormattedTeacherPhoneNumber string      `json:"FormattedTeacherPhoneNumber"`
		GuidanceApprovalDate        string      `json:"GuidanceApprovalDate"`
		GuardianApprovalDate        string      `json:"GuardianApprovalDate"`
		WillBeTakenInLab            interface{} `json:"WillBeTakenInLab"`
		ClassroomReservationID      interface{} `json:"ClassroomReservationId"`
		VirtualProgramID            int         `json:"VirtualProgramId"`
		EnrollmentID                int         `json:"EnrollmentId"`
		UserID                      int         `json:"UserId"`
		EnrollmentStatusID          int         `json:"EnrollmentStatusId"`
		CourseSegmentID             int         `json:"CourseSegmentId"`
		EnrollmentStatusDate        string      `json:"EnrollmentStatusDate"`
		VirtualSchoolID             int         `json:"VirtualSchoolId"`
		AverageGrade                float64     `json:"AverageGrade"`
		FinalGrade                  float64     `json:"FinalGrade"`
		FinalLetterGrade            string      `json:"FinalLetterGrade"`
		RowVersion                  string      `json:"RowVersion"`
		EnrollmentStatus            struct {
			EnrollmentStatusID      int    `json:"EnrollmentStatusId"`
			Description             string `json:"Description"`
			Abbreviation            string `json:"Abbreviation"`
			Stage                   int    `json:"Stage"`
			SortOrder               int    `json:"SortOrder"`
			ParentID                int    `json:"ParentId"`
			IsAlive                 bool   `json:"IsAlive"`
			IsStudentPlacementModel bool   `json:"IsStudentPlacementModel"`
			Group                   int    `json:"Group"`
			IsPending               bool   `json:"IsPending"`
			GroupName               string `json:"GroupName"`
		} `json:"EnrollmentStatus"`
	} `json:"Data"`
	Message string `json:"Message"`
}

// Class is the struct that holds class-specific information like grades
type Class struct {
	Grade           string
	PercentComplete float64
	LastSubmitted   string
	CourseName      string
}

// GetClasses - Gets the classes for a user
// First argument is the users id
// Second argument is the api key obtained from GetAPIKey function
func GetClasses(user, apiKey string) (*CourseInfo, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://vsaapi.flvs.net/api/students/%v/enrollments", user), nil)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	req.Header.Set("api-logon-token", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c := &CourseInfo{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetAPIKey - gets the main api key, consists of 10 functions gives back api key
func GetAPIKey(username, password string, verboose bool) (string, error) {
	if verboose {
		fmt.Println("Requesting to https://login.flvs.net/ with a post request.")
	}
	aspCookie1, nextLocation, err := login(username, password)
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a get request.\n", nextLocation)
	}
	aspxAuth1, aspCookie2, nextLocation, err := secondStep(nextLocation)
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to https://sts.flvs.net%v with a get request.\n", nextLocation)
	}
	postData, rptsSiteCookie1, nextLocation, err := thirdStep("https://sts.flvs.net"+nextLocation, fmt.Sprintf("ASP.NET_SessionId=%v; .ASPXAUTH=%v", aspCookie2, aspxAuth1))
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a POST request.\n", nextLocation)
	}
	_, err = fourthStep(nextLocation, fmt.Sprintf("ASP.NET_SessionId=%v; rememberMe=False", aspCookie1), postData)
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a get request.\n", "https://vsa.flvs.net/")
	}
	nextLocation, err = fifthStep()
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a get request.\n", nextLocation)
	}
	postData, _, nextLocation, err = sixthStep(nextLocation, fmt.Sprintf("ASP.NET_SessionId=%v; .ASPXAUTH=%v; RPStsSiteCookie=RPStsSite=%v", aspCookie2, aspxAuth1, rptsSiteCookie1))
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a post request.\n", nextLocation)
	}
	fedAuth2, nextLocation, err := seventhStep(postData, nextLocation)
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %s with a GET request.\n", nextLocation)
	}
	aspCookie3, nextLocation, err := eigthStep(nextLocation, fmt.Sprintf("FedAuth=%v", fedAuth2))
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a get request.\n", nextLocation)
	}
	svSID, setSEC, nextLocation, err := ninethStep(nextLocation, fmt.Sprintf("FedAuth=%v; ASP.NET_SessionId=%v", fedAuth2, aspCookie3))
	if err != nil {
		return "", err
	}
	if verboose {
		fmt.Printf("Requesting to %v with a get request.\n", nextLocation)
	}
	apiKey, err := tenthStep(nextLocation, fmt.Sprintf("ASP.NET_SessionId=%v; FedAuth=%v; SVSID=%v; %v", aspCookie3, fedAuth2, svSID, setSEC))
	if err != nil {
		return "", err
	}
	apiKey = parseAPIKey(apiKey)
	if verboose {
		fmt.Println("API key:", apiKey)
	}
	return apiKey, nil
}
