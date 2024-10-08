package plant

import "time"

type Plant struct {
	Name            string    `json:"name"`
	TimeLastWatered time.Time `json:"time_last_watered"`
}

//last watered method. Since only reading, we can just use a copy of Plant
func (p Plant) LastWatered() (days int64, ok bool) {
	if p.TimeLastWatered.IsZero() {
		return 0, false
	}
	day := 24 * time.Hour
	//only care about day(s). Truncate discards any non-day data.
	today := time.Now().Truncate(day)
	lastWatered := p.TimeLastWatered.Truncate(day)
	//.Sub -->> difference between today and lastWatered. Divide by day to
	//re-establish day as the units
	return int64(today.Sub(lastWatered) / day), true
}

//directly updating a value for plant struct. Must use a pointer
func (p *Plant) WaterMe() {
	p.TimeLastWatered = time.Now()
}
