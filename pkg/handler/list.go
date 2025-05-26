package handler

import (
	todo "do-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetAllListResponse struct {
	Data []todo.TodoList `json:"data"`
}

// @Summary Create list
// @Tags lists
// @Security ApiKeyAuth
// @Description create new list
// @ID new-list
// @Accept json
// @Produce json
// @Param input body todo.TodoList true "list information"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.TodoLists.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get All lists
// @Tags lists
// @Security ApiKeyAuth
// @Description get all lists
// @ID all-lists
// @Accept json
// @Produce json
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	lists, err := h.services.TodoLists.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, GetAllListResponse{
		Data: lists,
	})
}

// @Summary Get List By Id
// @Tags lists
// @Security ApiKeyAuth
// @Description get one list by id
// @ID id-lists
// @Accept json
// @Produce json
// @Param id path string true "list id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoLists.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary Update list
// @Tags lists
// @Security ApiKeyAuth
// @Description update list information
// @ID update-list
// @Accept json
// @Produce json
// @Param id path string true "list id"
// @Param input body todo.UpdateListInput false "information for update"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input todo.UpdateListInput
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err = h.services.TodoLists.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Delete list
// @Tags lists
// @Security ApiKeyAuth
// @Description delete list by id
// @ID delete-lists
// @Accept json
// @Produce json
// @Param id path string true "list id"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoLists.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
