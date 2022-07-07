package utils

type UserData struct {
	MsgId       string         `bson:"msgid,omitempty"`
	Name        string         `bson:"name,omitempty"`
	Age         string         `bson:"age,omitempty"`
	PhoneNo     string         `bson:"phone_no,omitempty"`
	Email       string         `bson:"email,omitempty"`
	ProfilePic  string         `bson:"profile_pic,omitempty"`
	MainKey     string         `bson:"main_key,omitempty"`
	Gender      string         `bson:"gender,omitempty"`
	Password    string         `bson:"password,omitempty"`
	Blocked     map[string]int `bson.M:"blocked,omitempty"`
	Connections map[string]int `bson.A:"connections,omitempty"`
}

type NewUser struct {
	MsgId       string                       `json:"id"`
	Name        string                       `json:"name"`
	Age         string                       `json:"age"`
	PhoneNo     string                       `json:"phoneno"`
	Email       string                       `json:"email"`
	ProfilePic  string                       `json:"profilepic"`
	MainKey     string                       `json:"mainkey"`
	Gender      string                       `json:"gender"`
	Password    string                       `json:"password"`
	Connections map[string]map[string]string `json:"connections"`
}

type MobileNo struct {
	Number string `json:"number"`
}

type VOTP struct {
	Otp string `json:"otp"`
	Id  string `json:"id"`
}

type SuccessStruct struct {
	Status string `json:"status"`
	Disc   string `json:"disc"`
}

type OperationId struct {
	Id string `json:"id"`
}

type LoginStruct struct {
	Target     string `json:"number"`
	Password   string `json:"password"`
	TargetType string `json:"target_type"`
}

type LoginSuccessPaylaod struct {
	MsgId      string `json:"id"`
	Name       string `json:"name"`
	Age        string `json:"age"`
	PhoneNo    string `json:"phone_no"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
	MainKey    string `json:"main_key"`
}

type ToggleBlocking struct {
	Type      int    `json:"type"`
	SenderMID string `json:"sendermid"`
	TargetNUM string `json:"number"`
}

// type UserData struct {
// 	IsOnline int
// 	LastSeen string
// 	Server   int // for getting the server name to which client is connected
// 	AuthOTP  int
// }

type DeletePayload struct {
	Target     string `json:"target"`
	TargetType string `json:"targettype"`
	OpId       string `json:"opid"`
}

type CheckAccPayload struct {
	SenderMID string `json:"sendermid"`
	TargetNUM string `json:"targetnum"`
}

type HandshakeDeletePaylaod struct {
	UserMID   string `json:"usermid"`
	TargetNUM string `json:"targetnum"`
}

// type HandshakeDeleteNotify struct {
// 	senderMID string `json:""`
// }

type ChangeProfilePayload struct {
	PicData   string `json:"picdata"`
	SenderMID string `json:"sendermid"`
}

type UpdateNumber struct {
	Otp    string `json:"otp"`
	Number string `json:"number"`
	MID    string `json:"mid"`
	OtpId  string `json:"otpid"`
	Notify int    `json:"notify"`
}

type UpdateEmailPayload struct {
	Email string `json:"email"`
	MID   string `json:"mid"`
}
