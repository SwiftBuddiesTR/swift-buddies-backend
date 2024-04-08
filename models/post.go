package post

type PostModel struct {
	ID      string `bson:"_id" json:"id"`
	Content string
}

type PostInput struct {
	Content  string `json:"content"`
	ImageUrl string `json:"imageUrl"`
	Text     string `json:"text"`
	UserId   string `json:"userId"`
}

type AddCategoryInput struct {
	Name string `json:"name" bson:"name"`
}

type AddItemInput struct {
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	AllComments []*string `json:"allComments,omitempty" bson:"allComments"`
	Category    string    `json:"category" bson:"category"`
}

type AddUserInput struct {
	Name        string    `json:"name" bson:"name"`
	Surname     string    `json:"surname" bson:"surname"`
	Email       string    `json:"email" bson:"email"`
	AllComments []*string `json:"allComments,omitempty" bson:"allComments"`
	IsVerified  *bool     `json:"isVerified,omitempty" bson:"isVerified"`
}

type Category struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type Comment struct {
	ID      string `json:"id" bson:"_id"`
	Content string `json:"content" bson:"content"`
	UserID  string `json:"userId" bson:"userId"`
	ItemID  string `json:"itemId" bson:"itemId"`
}

type Item struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Category    *Category `json:"category" bson:"category"`
}

type UpdateCategoryInput struct {
	Name *string `json:"name,omitempty" bson:"name"`
}

type UpdateItemInput struct {
	Name        *string   `json:"name,omitempty" bson:"name"`
	Description *string   `json:"description,omitempty" bson:"description"`
	AllComments []*string `json:"allComments,omitempty" bson:"allComments"`
	Category    *string   `json:"category,omitempty" bson:"category"`
}

type UpdateUserInput struct {
	Name        *string   `json:"name,omitempty" bson:"name"`
	Surname     *string   `json:"surname,omitempty" bson:"surname"`
	Email       *string   `json:"email,omitempty" bson:"email"`
	AllComments []*string `json:"allComments,omitempty" bson:"allComments"`
	IsVerified  *bool     `json:"isVerified,omitempty" bson:"isVerified"`
}

type User struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Surname     string    `json:"surname" bson:"surname"`
	Email       string    `json:"email" bson:"email"`
	AllComments []*string `json:"allComments,omitempty" bson:"allComments"`
	IsVerified  *bool     `json:"isVerified,omitempty" bson:"isVerified"`
}
