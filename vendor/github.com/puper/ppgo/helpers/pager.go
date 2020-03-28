package helpers

type PaginationArg struct {
	Page int
	Size int
}

type PaginationReply struct {
	Empty  bool
	Limit  int
	Offset int
	More   bool
}

func Pagination(arg *PaginationArg, totalSize, defaultPageSize, maxPageSize int) *PaginationReply {
	reply := new(PaginationReply)
	if arg.Page < 1 {
		arg.Page = 1
	}
	if arg.Size < 1 {
		arg.Size = defaultPageSize
	}
	if arg.Size > maxPageSize {
		arg.Size = maxPageSize
	}
	reply.Limit = arg.Size
	reply.Offset = arg.Size*arg.Page - arg.Size
	if totalSize < 0 {
		reply.More = true
	} else if totalSize == 0 {
		reply.Empty = true
	} else {
		if reply.Offset >= totalSize {
			reply.Empty = true
		} else if reply.Offset+arg.Size >= totalSize {
			reply.Limit = totalSize - reply.Offset
		} else {
			reply.More = true
		}
	}
	return reply
}
