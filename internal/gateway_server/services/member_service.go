package services

import (
	"poplargrid/internal/gateway_server/repositories"
	"poplargrid/internal/shared/dbmodels"
)

// MemberService 提供成员相关的服务
type MemberService struct {
	membersRepo *repositories.MembersRepo // 假设有一个成员仓库
}

// NewMemberService 创建一个新 MemberService
func NewMemberService(
	membersRepo *repositories.MembersRepo,
) *MemberService {
	return &MemberService{
		membersRepo: membersRepo,
	}
}

// GetMemberInfoById 获取成员信息
// 根据成员 ID 获取成员的完整信息
func (s *MemberService) GetMemberInfoById(memberId uint) (*dbmodels.User, error) {
	return s.membersRepo.SelectFullModelById(memberId)
}

// GetMemberInfoByNickname 获取成员信息
// 根据成员昵称获取成员的完整信息
func (s *MemberService) GetMemberInfoByNickname(nickname string) (*dbmodels.User, error) {
	return s.membersRepo.SelectFullModelByNickname(nickname)
}

// GetAllMembers 获取所有成员信息
// 返回所有成员的完整信息
func (s *MemberService) GetAllMembers() ([]dbmodels.User, error) {
	return s.membersRepo.SelectAllToSlice()
}
