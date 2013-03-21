package hpcloud

/*
 BadRequest describes the response from a JSON resource when the
 data which was sent in the original request was malformed or not
 compliant with the layout specified in the HPCloud documentation
*/
type BadRequest struct {
	B struct {
		Message string `json:"message"`
		Details string `json:"details"`
		Code    int64  `json:"code"`
	} `json:"BadRequest"`
}

/*
 NotFound describes the response from a JSON resource when the
 resource which was interacted with in the original request was
 not able to be found.
*/
type NotFound struct {
	NF struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"itemNotFound"`
}

/*
 Unauthorized describes the response from a JSON resource when the
 request could not be completed due to none or incorrect authentication
 was used to make the request.
*/
type Unauthorized struct {
	U struct {
		Code            int64  `json:"code"`
		Details         string `json:"details"`
		Message         string `json:"message"`
		OtherAttributes struct {
		} `json:"otherAttributes"`
	} `json:"unauthorized"`
}

type SubToken struct {
	ID string `json:"id"`
}

type Tenants struct {
	T []Tenant `json:"tenants"`
}

type FailureResponse interface {
	Code() int64
	Details() string
	Message() string
}

func (u Unauthorized) Code() int64 {
	return u.U.Code
}

func (b BadRequest) Code() int64 {
	return b.B.Code
}

func (u Unauthorized) Details() string {
	return u.U.Details
}

func (b BadRequest) Details() string {
	return b.B.Details
}

func (u Unauthorized) Message() string {
	return u.U.Message
}

func (b BadRequest) Message() string {
	return b.B.Message
}
