package handler

type Handler struct {
	MsgHandler   *MsgHandler
	UserHandler  *UserHandler
	PostHandler  *PostHandler
	CacheHandler *CacheHandler
	QueueHandler *QueueHandler
	ImageHandler *ImageHandler
}
