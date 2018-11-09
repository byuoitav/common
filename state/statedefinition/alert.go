package statedefinition

import "time"

type Alert struct {
	AlertSent time.Time `json:"alert-sent,omitempty"`
	Alerting  bool      `json:"alerting,omitempty"`
	Message   string    `json:"message,omitempty"`
}

func compareAlerts(base, new map[string]Alert, basetime, newtime map[string]time.Time, changes bool) (diff map[string]Alert, merged map[string]Alert, c bool) {
	for k, v := range new {

		basev, ok := base[k]
		if !ok {
			c = true
			base[k] = v
			diff[k] = v
		}

		if newtime["alerts."+k].Before(basetime["alerts."+k]) {
			continue
		}

		new, tempChanges := compareAlert(basev, v)
		if tempChanges {
			c = true
			base[k] = new
			diff[k] = new
		}
	}
	return
}

func compareAlert(base, new Alert) (after Alert, changes bool) {
	if base.Alerting != new.Alerting || base.AlertSent.Equal(new.AlertSent) || base.Message != new.Message {
		after = new
		changes = true
		return
	}

	return
}
