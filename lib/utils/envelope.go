package utils

import (
	"net/http"
	"time"
)

type Envelope struct {
	Success        bool        `json:"success"`
	RequestID      string      `json:"request_id"`
	Time           string      `json:"time"`
	Response       interface{} `json:"response"`
	Error          string      `json:"error,omitempty"`
	Message        string      `json:"message,omitempty"`
	HTTPStatusCode int         `json:"-"`
}

func NewEnvelope() *Envelope {
	return &Envelope{}
}

func (e *Envelope) Render(w http.ResponseWriter, r *http.Request) error {
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
