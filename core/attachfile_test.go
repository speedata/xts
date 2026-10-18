package core

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A Factur-X invoice with the profile MINIMUM or BASIC WL is attached with
// /AFRelationship /Data, one with a richer profile with /Alternative.
func TestAttachFileFacturxAFRelationship(t *testing.T) {
	for guideline, want := range map[string]string{
		"urn:factur-x.eu:1p0:minimum": "Data",
		"urn:factur-x.eu:1p0:basicwl": "Data",
		"urn:cen.eu:en16931:2017":     "Alternative",
	} {
		dir := t.TempDir()
		invoice := filepath.Join(dir, "factur-x.xml")
		xml := `<rsm:CrossIndustryInvoice xmlns:rsm="urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100" xmlns:ram="urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100">
  <rsm:ExchangedDocumentContext><ram:GuidelineSpecifiedDocumentContextParameter><ram:ID>` + guideline + `</ram:ID></ram:GuidelineSpecifiedDocumentContextParameter></rsm:ExchangedDocumentContext>
</rsm:CrossIndustryInvoice>`
		if err := os.WriteFile(invoice, []byte(xml), 0o644); err != nil {
			t.Fatal(err)
		}
		jobname := filepath.Join(dir, "out")
		out, err := os.Create(jobname + ".pdf")
		if err != nil {
			t.Fatal(err)
		}
		err = RunXTS(&XTSConfig{
			Datafile: strings.NewReader("<data/>"),
			Layoutfile: strings.NewReader(layoutHead + `
  <Record match="data">
    <AttachFile href="` + invoice + `" type="facturx"/>
    <PlaceObject><TextBlock><Paragraph><Value>Invoice</Value></Paragraph></TextBlock></PlaceObject>
  </Record>
</Layout>`),
			FindFile:    FindFile,
			Jobname:     jobname,
			Outfile:     out,
			OutFilename: out.Name(),
		})
		out.Close()
		if err != nil {
			t.Fatal(err)
		}
		pdf, err := os.ReadFile(jobname + ".pdf")
		if err != nil {
			t.Fatal(err)
		}
		m := regexp.MustCompile(`/AFRelationship /(\w+)`).FindSubmatch(pdf)
		if m == nil || string(m[1]) != want {
			t.Errorf("%s: AFRelationship %q, want %q", guideline, m, want)
		}
	}
}
