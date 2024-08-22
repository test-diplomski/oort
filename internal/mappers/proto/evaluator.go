package proto

import (
	"log"

	"github.com/c12s/oort/internal/domain"
	"github.com/c12s/oort/pkg/api"
)

func AuthorizationReqToDomain(req *api.AuthorizationReq) (*domain.AuthorizationReq, error) {
	envAttributes := make([]domain.Attribute, len(req.EnvAttributes))
	for i, attr := range req.EnvAttributes {
		domainAttr, err := AttributeToDomain(attr)
		if err != nil {
			log.Println(err)
			continue
		}
		envAttributes[i] = *domainAttr
	}
	sub, err := ResourceToDomain(req.Subject)
	if err != nil {
		return nil, err
	}
	obj, err := ResourceToDomain(req.Object)
	if err != nil {
		return nil, err
	}
	return &domain.AuthorizationReq{
		Subject:        *sub,
		Object:         *obj,
		PermissionName: req.PermissionName,
		Env:            envAttributes,
	}, nil
}

func GetGrantedPermissionsReqToDomain(req *api.GetGrantedPermissionsReq) (*domain.GetGrantedPermissionsReq, error) {
	envAttributes := make([]domain.Attribute, len(req.EnvAttributes))
	for i, attr := range req.EnvAttributes {
		domainAttr, err := AttributeToDomain(attr)
		if err != nil {
			log.Println(err)
			continue
		}
		envAttributes[i] = *domainAttr
	}
	sub, err := ResourceToDomain(req.Subject)
	if err != nil {
		return nil, err
	}
	return &domain.GetGrantedPermissionsReq{
		Subject: *sub,
		Env:     envAttributes,
	}, nil
}

func GetGrantedPermissionsRespFromDomain(resp *domain.GetGrantedPermissionsResp) (*api.GetGrantedPermissionsResp, error) {
	perms := make([]*api.GrantedPermission, 0)
	for _, domainPerm := range resp.Permissions {
		perm, err := GrantedPermissionFromDomain(&domainPerm)
		if err != nil {
			log.Println(err)
			continue
		}
		perms = append(perms, perm)
	}
	return &api.GetGrantedPermissionsResp{
		Permissions: perms,
	}, nil
}
