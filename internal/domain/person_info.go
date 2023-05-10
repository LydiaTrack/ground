package domain

type PersonInfo struct {
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Email     string `bson:"email,omitempty"`
	Address   string `bson:"address,omitempty"`
	Phone     string `bson:"phone,omitempty"`
}
