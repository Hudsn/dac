package dac

func (dr *DetectionResult) SetTitle(title string) {
	dr.Title = title
}

func (dr *DetectionResult) SetDescription(desc string) {
	dr.Description = desc
}

func (dr *DetectionResult) SetClassification(c Classification) {
	dr.Classification = c
}

func (dr *DetectionResult) SetMatch(tf bool) {
	dr.IsMatch = tf
}

func (dr *DetectionResult) AddNote(note DetectionNote) {
	dr.Notes = append(dr.Notes, note)
}

func (dr DetectionResult) HasTimestamp() bool {
	return !dr.Timestamp.IsZero()
}

func (dr DetectionResult) HasTitle() bool {
	return len(dr.Title) > 0
}

func (dr DetectionResult) HasDescription() bool {
	return len(dr.Description) > 0
}

func (dr DetectionResult) HasClassification() bool {
	return len(dr.Classification) > 0
}

func (dr DetectionResult) HasNotes() bool {
	return len(dr.Notes) > 0 || dr.Notes != nil
}
