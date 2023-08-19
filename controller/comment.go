package controller

import (
	"douyin/model"
	"douyin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	model.Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	var request model.CommentActionRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusOK, model.CommentActionResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	user, exist := service.CheckLogin(request.Token)
	if !exist {
		c.JSON(http.StatusOK, model.CommentActionResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	if request.ActionType == 1 { // Post Comment
		commentDAO, err := service.CommentPost(user.Id, request.VideoId, request.Content)
		if err != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
			Comment: model.Comment{
				Id:         commentDAO.ID,
				User:       *user,
				Content:    commentDAO.Content,
				CreateDate: commentDAO.CreatedAt.Format("01-02"),
			}})
		// 测试数据
		/*
			c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
				Comment: model.Comment{
					Id:         1,
					User:       *user,
					Content:    text,
					CreateDate: "05-01",
				}})
		*/
		return
	}
	if request.ActionType == 2 { // Delete Comment
		err := service.CommentDelete(request.CommentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: model.Response{StatusCode: 0},
	})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	_, exist := service.CheckLogin(token)
	if !exist {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	commentListDAO, err := service.CommentList(videoId)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	var commentList []model.Comment = make([]model.Comment, len(commentListDAO))
	for i, commentDAO := range commentListDAO {
		commentList[i] = model.Comment{
			Id: commentDAO.ID,
			User: model.User{
				Id:            commentDAO.User.ID,
				Name:          commentDAO.User.Name,
				FollowCount:   0,     // TODO
				FollowerCount: 0,     // TODO
				IsFollow:      false, // TODO
			},
			Content:    commentDAO.Content,
			CreateDate: commentDAO.CreatedAt.Format("01-02"),
		}
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: commentList,
	})

	// 测试数据
	/*
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    model.Response{StatusCode: 0},
			CommentList: DemoComments,
		})
	*/
}
