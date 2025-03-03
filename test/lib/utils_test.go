package lib

import (
	"maps-to-waze-api/lib"
	"testing"
)

type convertLinkData struct {
	url      string
	expected string
}

func TestConvertLinkExpanded(t *testing.T) {
	testCases := []convertLinkData{
		{
			"https://www.google.com/maps/place/Scuola+Elementare+le+Risorgive/@45.3792597,10.9885404,727m/data=!3m2!1e3!4b1!4m6!3m5!1s0x477f5fd9d5d3fa77:0x6361b6fc9329ef1c!8m2!3d45.3792597!4d10.9885404!16s%2Fg%2F1tdc_jf4?entry=ttu&g_ep=EgoyMDI1MDIyNi4xIKXMDSoASAFQAw%3D%3D",
			"https://www.waze.com/ul?ll=45.3792597,10.9885404&navigate=yes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			testConvertLink(t, tc.url, tc.expected)
		})
	}
}

func TestConvertLinkCompressed(t *testing.T) {
	testCases := []convertLinkData{
		{
			"https://maps.app.goo.gl/esEFcSdtrwfCs7es9",
			"https://www.waze.com/ul?ll=45.5571256,10.7701129&navigate=yes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			testConvertLink(t, tc.url, tc.expected)
		})
	}
}

func TestConvertLinkWithApi(t *testing.T) {
	testCases := []convertLinkData{
		{
			"https://maps.app.goo.gl/a6fd6rmhQxkRZvSc6?g_st=it",
			"https://www.waze.com/ul?ll=45.4258881,11.0723848&navigate=yes",
		},
		{
			"https://maps.app.goo.gl/GophJ1CwCUMW2YV97",
			"https://www.waze.com/ul?ll=45.3792597,10.9885404&navigate=yes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			testConvertLink(t, tc.url, tc.expected)
		})
	}
}

func testConvertLink(t *testing.T, url string, expected string) {
	wazeUrl, err := lib.ConvertLink(url)

	// Check if an error occurred
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// Verify that the result is not empty or invalid
	if wazeUrl == "" {
		t.Errorf("Expected a valid Waze URL, but got an empty string")
		return
	}

	// Check that the output is like expected
	if wazeUrl != expected {
		t.Errorf("Expected %s, but got: %s", expected, wazeUrl)
	}
}
