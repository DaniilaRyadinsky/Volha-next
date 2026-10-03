package health

type Status struct {
	Status string `json:"status"`
}

func Live() Status {
	return Status{Status: "ok"}
}

func Ready(err error) (Status, int) {
	if err != nil {
		return Status{Status: "not ready"}, 503
	}
	return Status{Status: "ready"}, 200
}
