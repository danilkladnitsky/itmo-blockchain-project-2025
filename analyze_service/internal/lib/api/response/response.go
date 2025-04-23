package response

type Response struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"OK"`
}

type ErrorResponse struct {
	Code    int   `json:"code" example:"400"`
	Message error `json:"message" example:"Bad Request"`
}

const (
	StatusOK   = 200
	StatusFail = 400
	StatusErr  = 500
)

func OK() Response {
	return Response{
		Code: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Code:    StatusFail,
		Message: msg,
	}
}
func ErrResponse(msg error) ErrorResponse {
	return ErrorResponse{
		400,
		msg,
	}
}

//func ResponseOK(w http.ResponseWriter, r *http.Request, model any) {
//	render.JSON(w, r, model)
//}

//func RenderErrorResponse(err error) {
//
//	return render.JSON(w, r, resp.Error("empty request"))
//}

//func ValidationError(errs validator.ValidationErrors) Response {
//	var errMsgs []string
//
//	for _, err := range errs {
//		switch err.ActualTag() {
//		case "required":
//			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
//		case "url":
//			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
//		default:
//			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
//		}
//	}
//
//	return Response{
//		Status: StatusError,
//		Error:  strings.Join(errMsgs, ", "),
//	}
//}
