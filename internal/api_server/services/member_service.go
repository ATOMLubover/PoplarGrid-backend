package services

import (
	"poplargrid/internal/api_server/repository"
	"poplargrid/internal/shared/dbmodel"
)

// MemberService 提供成员相关的服务
type MemberService struct {
	memberRepo *repository.MembersRepo
}

// NewMemberService 创建一个新的 MemberService 实例
func NewMemberService(
	memberRepo *repository.MembersRepo,
) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
	}
}

// GetMemberFullList 获取所有成员的完整信息
func (s *MemberService) GetMemberFullList(offset, num int) ([]*dbmodel.Member, error) {
	return s.memberRepo.SelectAllToSlice(offset, num)
}
