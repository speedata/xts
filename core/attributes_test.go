package core

import (
	"strings"
	"testing"
)

// An unknown attribute is reported once with its line and command, even when
// the command runs several times, with the attribute that was probably meant.
func TestUnknownAttributeWarnsOnce(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <Loop select="3">
      <PlaceObject colum="3"><Box width="2" height="1"/></PlaceObject>
    </Loop>
  </Record>
</Layout>`)
	want := `Layout line 4: unknown attribute \"colum\" on PlaceObject, did you mean \"column\"?`
	if n := strings.Count(log, want); n != 1 {
		t.Errorf("warning %s logged %d times, want once, in\n%s", want, n, log)
	}
}

// A name far from every attribute of the command gets no suggestion.
func TestUnknownAttributeWithoutSuggestion(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <PlaceObject><Box width="2" height="1" colro="red"/></PlaceObject>
  </Record>
</Layout>`)
	want := `Layout line 3: unknown attribute \"colro\" on Box"`
	if !strings.Contains(log, want) {
		t.Errorf("no warning %s in\n%s", want, log)
	}
}

// Attribute names are matched exactly: the hyphen of background-color is
// required, and the spelling without it is unknown.
func TestHyphenatedAttributeNames(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <PlaceObject><Box width="2" height="1" background-color="red"/></PlaceObject>
    <PlaceObject><Box width="2" height="1" backgroundcolor="red"/></PlaceObject>
  </Record>
</Layout>`)
	want := `Layout line 4: unknown attribute \"backgroundcolor\" on Box, did you mean \"background-color\"?`
	if !strings.Contains(log, want) {
		t.Errorf("no warning %s in\n%s", want, log)
	}
	if n := strings.Count(log, "unknown attribute"); n != 1 {
		t.Errorf("%d warnings about unknown attributes, want 1, in\n%s", n, log)
	}
}

// Commands without attributes, and the commands that their parent reads,
// report all attributes as unknown.
func TestUnknownAttributeOnCommandsWithoutAttributes(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <PlaceObject>
      <TextBlock>
        <Paragraph><B weight="bold">Hello</B></Paragraph>
      </TextBlock>
    </PlaceObject>
    <Switch>
      <Case test="false()" tset="true()"/>
      <Otherwise select="1"/>
    </Switch>
  </Record>
</Layout>`)
	for _, want := range []string{
		`Layout line 5: unknown attribute \"weight\" on B"`,
		// test is there, so tset is not taken for it.
		`Layout line 9: unknown attribute \"tset\" on Case"`,
		`Layout line 10: unknown attribute \"select\" on Otherwise"`,
	} {
		if !strings.Contains(log, want) {
			t.Errorf("no warning %s in\n%s", want, log)
		}
	}
}

// Attributes in a namespace are not checked, and a layout that uses only
// known attributes, including those that nothing but the schema documents,
// logs no warning.
func TestNoWarningForKnownAttributes(t *testing.T) {
	log := runLayoutLog(t, `<Layout xmlns="urn:speedata.de/2021/xts/en" xmlns:sd="urn:speedata.de/2021/xtsfunctions/en" xmlns:x="urn:example" version="0.1">
  <DefineMasterPage name="page" test="true()" margin="1cm">
    <AtPageCreation/>
    <PositioningArea name="text">
      <PositioningFrame width="10" height="10" row="1" column="1"/>
    </PositioningArea>
  </DefineMasterPage>
  <Record match="data">
    <Section name="intro" x:note="1">
      <PlaceObject xml:lang="en"><Box width="2" height="1" background-color="red"/></PlaceObject>
    </Section>
    <Switch>
      <Case test="false()"/>
      <Otherwise/>
    </Switch>
  </Record>
</Layout>`)
	if strings.Contains(log, "unknown attribute") {
		t.Errorf("warning about a known attribute in\n%s", log)
	}
}

func TestClosestName(t *testing.T) {
	names := []string{"column", "row", "area", "background-color"}
	for _, tc := range []struct{ name, want string }{
		{"colum", "column"},
		{"rwo", "row"},
		{"backgroundcolor", "background-color"},
		{"backgrund-colr", "background-color"},
		{"colro", ""},
		{"x", ""},
	} {
		if got := closestName(tc.name, names); got != tc.want {
			t.Errorf("closestName(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// append and shiftup of Mark were documented but had no effect; they are
// unknown now, so a layout that sets them hears about it.
func TestMarkAppendAndShiftupAreUnknown(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <Mark select="'m'" append="yes" shiftup="2mm"/>
    <PlaceObject><Box width="2" height="1"/></PlaceObject>
  </Record>
</Layout>`)
	for _, att := range []string{"append", "shiftup"} {
		if want := `unknown attribute \"` + att + `\" on Mark`; !strings.Contains(log, want) {
			t.Errorf("no warning %s in\n%s", want, log)
		}
	}
}
