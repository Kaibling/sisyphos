package utils

import (
	"net/http"
	"sisyphos/lib/reqctx"
	"time"
)

type LogService interface {
	Log(url, body, method, user, requestID string) error
}

type Envelope struct {
	Success        bool        `json:"success"`
	RequestID      string      `json:"request_id"`
	Time           string      `json:"time"`
	Response       interface{} `json:"response"`
	Error          string      `json:"error,omitempty"`
	Message        string      `json:"message,omitempty"`
	HTTPStatusCode int         `json:"-"`
	LogService     LogService  `json:"-"`
}

func NewEnvelope(ls LogService) *Envelope {
	return &Envelope{LogService: ls}
}

func (e *Envelope) Render(w http.ResponseWriter, r *http.Request) error {
	url := r.URL.String()
	body := string(reqctx.GetContext("bytebody", r).([]uint8))
	method := r.Method
	var user string
	if u, ok := reqctx.GetContext("username", r).(string); ok {
		user = u
	} else {
		user = "unauthenticated"
	}
	requestID := reqctx.GetContext("requestid", r).(string)
	e.LogService.Log(url, body, method, user, requestID)
	return nil
}

func (e *Envelope) SetResponse(resp interface{}) *Envelope {
	e.Success = true
	e.Time = time.Now().Format(time.RFC822)
	e.Response = resp
	// FillNull(resp)
	return e
}

func (e *Envelope) SetError(resp error) *Envelope {
	e.Success = false
	e.Time = time.Now().Format(time.RFC822)
	e.Response = resp.Error()
	// FillNull(resp)
	return e
}

// func FillNull(d interface{}) interface{} {

// 	v := reflect.ValueOf(d)
// 	typeOfS := v.Type()
// 	//v.Len()
// 	switch typeOfS.Kind() {
// 	case reflect.Slice | reflect.Array:
// 		for i := range d.([]interface{}) {
// 			//d[i]
// 			vara := d.([]interface{})[i]
// 		}
// 		fmt.Println("sssss")
// 	//case reflect.Array:
// 	default:
// 		fmt.Println("ss")
// 	}

// 	for i := 0; i < v.NumField(); i++ {
// 		//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
// 		fmt.Printf("Field: %s\tValue: \n", typeOfS.Field(i).Name)
// 	}
// 	return d
// }
