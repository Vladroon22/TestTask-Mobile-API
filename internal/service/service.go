package service

type User struct {
	Nick     string `json:"nick"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	Author  string `json:"poster"`
	Title   string `json:"title"`
	Content string `json:"text"`
}

type Comment struct {
	Author  string `json:"commenter"`
	Content string `json:"comment"`
}

type Like struct {
	LikerNick      string `json:"liker"`
	LikeOrDis      string `json:"type"`
	PostAuthor     string `json:"post_author"`
	PostContent    string `json:"post_content"`
	CommentAuthor  string `json:"comment_author"`
	CommentContent string `json:"comment_content"`
}

type Service struct {
	User
	Post
	Comment
	Like
}

func NewService() *Service {
	return &Service{}
}
