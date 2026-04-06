package cognitive

import "errors"

// ClassifyError matches the README's ResponseFromError example pattern.
// switch(n=0)=+1; each case(flat)=+1; branch inside each case(n=1)=+2
func ClassifyError(code int) error {
	switch code { // n=0 → +1
	case 400:
		return errors.New("bad request") // n=1 → +2
	case 401:
		return errors.New("unauthorized") // n=1 → +2
	case 403:
		return errors.New("forbidden") // n=1 → +2
	case 404:
		return errors.New("not found") // n=1 → +2
	case 500:
		return errors.New("internal error") // n=1 → +2
	default:
		return errors.New("unknown") // n=1 → +2
	}
}

// TypeDispatch shows a type switch.
func TypeDispatch(v interface{}) string {
	switch t := v.(type) { // n=0 → +1
	case int:
		return "int" // n=1 → +2
	case string:
		if len(t) == 0 { // n=1 → +2
			return "empty string" // n=2 → +3
		}
		return "string" // n=2 → +3
	default:
		return "other" // n=1 → +2
	}
}
