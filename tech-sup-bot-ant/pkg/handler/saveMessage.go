package handler

/*func (h *Handler) createMessage(c *gin.Context) {
	var input *domain.Message

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.MessageCRUD.CreateMessage(input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllMessage(c *gin.Context) {

	lists, err := h.services.MessageCRUD.GetAllMessage()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"messages": lists,
	})
}*/
