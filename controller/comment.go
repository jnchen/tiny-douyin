package controller

import (
	"douyin/model"
	"douyin/service"
	"douyin/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommentAction 发布评论或删除评论
func CommentAction(c *gin.Context) {
	var request model.CommentActionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, model.CommentActionResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	user := util.GetUser(c)
	if request.ActionType == 1 { // Post Comment
		commentDAO, err := service.CommentPost(user.Id, request.VideoId, request.Content)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				model.Response{StatusCode: 1, StatusMsg: err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, model.CommentActionResponse{
			Response: model.Response{StatusCode: 0},
			Comment:  *commentDAO.ToModel(),
		})
		return
	}
	if request.ActionType == 2 { // Delete Comment
		err := service.CommentDelete(request.CommentId, request.VideoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CommentActionResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, model.CommentActionResponse{
		Response: model.Response{StatusCode: 0},
	})
}

// CommentList 获取视频评论列表
func CommentList(c *gin.Context) {
	var req model.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.CommentListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	commentListDAO, err := service.CommentList(req.VideoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.CommentListResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	var commentList = make([]model.Comment, len(commentListDAO))
	for i, commentDAO := range commentListDAO {
		commentList[i] = *commentDAO.ToModel()
	}

	c.JSON(http.StatusOK, model.CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: commentList,
	})
}
