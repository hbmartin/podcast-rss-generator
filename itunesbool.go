package podcast

import (
	"encoding/xml"
	"fmt"
)

// explicitFlag renders the itunes:explicit parental-advisory tag, whose value
// is the literal "true" or "false" per Apple's specification. Unlike yesFlag,
// a false value is meaningful (it renders the "Clean" badge), so it is
// serialized rather than omitted.
type explicitFlag bool

// MarshalXML writes "true" or "false" for the itunes:explicit tag.
func (f explicitFlag) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	value := "false"
	if f {
		value = "true"
	}
	if err := e.EncodeElement(value, start); err != nil {
		return fmt.Errorf("encoding itunes:explicit: %w", err)
	}
	return nil
}

// yesFlag renders itunes tags such as complete and block. Apple only
// recognizes the literal value "Yes"; any other value has no effect, so a
// false flag is represented by omitting the element entirely.
type yesFlag bool

// MarshalXML writes "Yes" when the flag is true and nothing when it is false.
func (f yesFlag) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !f {
		return nil
	}
	if err := e.EncodeElement("Yes", start); err != nil {
		return fmt.Errorf("encoding %s: %w", start.Name.Local, err)
	}
	return nil
}
