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
		Code    int64  `json:"code"`
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

/*
 Forbidden describes the response from a JSON resource when the
 request could not be completed due to the user making the request
 being disabled or suspended.
*/
type Forbidden struct {
	F struct {
		Code            int64  `json:"code"`
		Details         string `json:"details"`
		Message         string `json:"message"`
		OtherAttributes struct {
		} `json:"otherAttributes"`
	} `json:"forbidden"`
}

/*
 InternalServerError describes the response from a JSON resource when the
 request could not be completed due to the request causing the service
 to return a 500 status code.
*/
type InternalServerError struct {
	ISE struct {
		Code            int64  `json:"code"`
		Details         string `json:"details"`
		Message         string `json:"message"`
		OtherAttributes struct {
		} `json:"otherAttributes"`
	} `json:"internalServerError"`
}

type SubToken struct {
	ID string `json:"id"`
}

type Flavor_ struct {
	Name  string `json:"name"`
	ID    int64  `json:"id"`
	Links []Link `json:"links"`
}

type Flavors struct {
	F []Flavor_ `json:"flavors"`
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

func (f Forbidden) Code() int64 {
	return f.F.Code
}

func (ise InternalServerError) Code() int64 {
	return ise.ISE.Code
}

func (nf NotFound) Code() int64 {
	return nf.NF.Code
}

func (u Unauthorized) Details() string {
	return u.U.Details
}

func (b BadRequest) Details() string {
	return b.B.Details
}

func (f Forbidden) Details() string {
	return f.F.Details
}

func (ise InternalServerError) Details() string {
	return ise.ISE.Details
}

func (u Unauthorized) Message() string {
	return u.U.Message
}

func (b BadRequest) Message() string {
	return b.B.Message
}

func (f Forbidden) Message() string {
	return f.F.Message
}

func (ise InternalServerError) Message() string {
	return ise.ISE.Message
}

func (nf NotFound) Message() string {
	return nf.NF.Message
}
