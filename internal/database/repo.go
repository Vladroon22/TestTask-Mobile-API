package database

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/Vladroon22/TestTask-Mobile-API/internal/service"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/utils"
	"github.com/dgrijalva/jwt-go"
)

const (
	TTLofJWT = time.Minute * 10
)

var SignKey = []byte(os.Getenv("KEY"))

type MyClaims struct {
	jwt.StandardClaims
	UserID int
}

type Repo struct {
	db *DataBase
}

func NewRepo(db *DataBase) *Repo {
	return &Repo{db: db}
}

func (rp *Repo) CreateUser(user *service.User) error {
	enc_pass, err := utils.Hashing(user.Password)
	if err != nil {
		rp.db.logger.Errorln(err)
		return err
	}
	query := "INSERT INTO users (nickname, name, email, hash) VALUES ($1, $2, $3, $4)"
	if _, err := rp.db.sqlDB.Exec(query, user.Nick, user.Name, user.Email, string(enc_pass)); err != nil {
		rp.db.logger.Errorln(err)
		return err
	}

	rp.db.logger.Infoln("User successfully added")
	return nil
}

func (rp *Repo) Login(pass, email string) (int, error) {
	var id int
	var hash string

	query1 := "SELECT id, hash FROM users WHERE email = $1"
	if err := rp.db.sqlDB.QueryRow(query1, email).Scan(&id, &hash); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Wrong-password-or-email")
		}
		return 0, err
	}

	if err := utils.CheckPassAndHash(hash, pass); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	return id, nil
}

func (rp *Repo) Posting(post *service.Post) error {
	query := "INSERT INTO posts (author, title, content) VALUES ($1, $2, $3)"
	if _, err := rp.db.sqlDB.Exec(query, post.Author, post.Title, post.Content); err != nil {
		return err
	}
	return nil
}

func (rp *Repo) GetPost(id int) (*service.Post, error) {
	storedPost := &service.Post{}

	query := `
	SELECT 
		users.nickname AS user_nick, 
		posts.title AS post_title, 
		posts.content AS post_content
	FROM posts 
	LEFT JOIN 
		users ON posts.user_id = users.id
	WHERE posts.id = $1`
	if err := rp.db.sqlDB.QueryRow(query, id).Scan(&storedPost.Author, &storedPost.Title, &storedPost.Content); err != nil {
		return nil, err
	}
	return storedPost, nil
}

func (rp *Repo) Comment(comm *service.Comment) error {
	query := "INSERT INTO comments (author, content) VALUES ($1, $2)"
	if _, err := rp.db.sqlDB.Exec(query, comm.Author, comm.Content); err != nil {
		return err
	}
	return nil
}

func (rp *Repo) GetComment(id int) (*service.Comment, error) {
	storedComm := &service.Comment{}
	query := `
	SELECT 
		users.nickname AS user_nick, 
		comments.content AS comment_content 
	FROM comments 
	LEFT JOIN 
		users ON comments.user_id = users.id
	WHERE comments.id = $1`
	if err := rp.db.sqlDB.QueryRow(query, id).Scan(&storedComm.Author, &storedComm.Content); err != nil {
		return nil, err
	}
	return storedComm, nil
}

func (rp *Repo) LikeIt(nick, like, type_cont string, id int) error {
	var query string

	if type_cont == "post" {
		query = "INSERT INTO likes (type_of_like, liker, post_id) VALUES ($1, $2, $3)"
	} else if type_cont == "comment" {
		query = "INSERT INTO likes (type_of_like, liker, comment_id) VALUES ($1, $2, $3)"
	}

	if _, err := rp.db.sqlDB.Exec(query, like, nick, id); err != nil {
		return err
	}

	return nil
}

func (rp *Repo) GetLiker(id int) (*service.Like, error) {
	storedLike := &service.Like{}
	query := `
        SELECT 
            liker, 
            type_of_like, 
            posts.author AS post_author, 
            posts.content AS post_content, 
            comments.author AS comment_author, 
            comments.content AS comment_content 
        FROM 
            likes 
        LEFT JOIN 
            comments ON likes.comment_id = comments.id 
        LEFT JOIN 
            posts ON likes.post_id = posts.id 
        WHERE 
            likes.id = $1`

	if err := rp.db.sqlDB.QueryRow(query, id).
		Scan(
			&storedLike.LikerNick,
			&storedLike.LikeOrDis,
			&storedLike.PostAuthor,
			&storedLike.PostContent,
			&storedLike.CommentAuthor,
			&storedLike.CommentContent); err != nil {
		return nil, err
	}
	return storedLike, nil
}

func (rp *Repo) GenerateJWT(id int) (string, error) {
	JWT, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TTLofJWT).Unix(), // TTL of token
			IssuedAt:  time.Now().Unix(),
		},
		id,
	}).SignedString(SignKey)
	if err != nil {
		return "", err
	}

	return JWT, nil
}

func ValidateToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SignKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New(err.Error())
		}
		return nil, errors.New(err.Error())
	}

	claims, ok := token.Claims.(*MyClaims)
	if !token.Valid {
		return nil, errors.New("Token-is-invalid")
	}
	if !ok {
		return nil, errors.New("Unauthorized")
	}

	return claims, nil
}
