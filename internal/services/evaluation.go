package services

import (
	"log"

	"github.com/c12s/oort/internal/domain"
)

type EvaluationService struct {
	repo domain.RHABACRepo
}

func NewEvaluationService(repo domain.RHABACRepo) (*EvaluationService, error) {
	return &EvaluationService{
		repo: repo,
	}, nil
}

func (h EvaluationService) Authorize(req domain.AuthorizationReq) domain.AuthorizationResp {
	resp := h.repo.GetPermissionHierarchy(domain.GetPermissionHierarchyReq{
		Subject:        req.Subject,
		Object:         req.Object,
		PermissionName: req.PermissionName,
	})
	if resp.Error != nil {
		return domain.AuthorizationResp{
			Authorized: false,
			Error:      resp.Error,
		}
	}

	subAttrs, err := h.getAttributes(req.Subject)
	if err != nil {
		return domain.AuthorizationResp{
			Authorized: false,
			Error:      err,
		}
	}
	objAttrs, err := h.getAttributes(req.Object)
	if err != nil {
		return domain.AuthorizationResp{
			Authorized: false,
			Error:      err,
		}
	}

	evalReq := domain.PermissionEvalRequest{
		Subject: subAttrs,
		Object:  objAttrs,
		Env:     req.Env,
	}
	evalResult := resp.Hierarchy.Eval(evalReq)

	checkResp := domain.AuthorizationResp{
		Authorized: authorized(evalResult),
		Error:      nil,
	}

	return checkResp
}

func (h EvaluationService) GetGrantedPermissions(req domain.GetGrantedPermissionsReq) domain.GetGrantedPermissionsResp {
	// dobavi sve politike koje su subjektno direktno dodeljene ili ih je nasledio
	// svaka ukljucuje naziv dozvole i objekat nad kojim vazi
	resp := h.repo.GetApplicablePolicies(domain.GetApplicablePoliciesReq{
		Subject: req.Subject,
	})
	if resp.Error != nil {
		return domain.GetGrantedPermissionsResp{Error: resp.Error}
	}

	granted := make([]domain.GrantedPermission, 0)

	subAttrs, err := h.getAttributes(req.Subject)
	if err != nil {
		return domain.GetGrantedPermissionsResp{Error: resp.Error}
	}
	// proveravamo nad vise objekata, svaki objekat je element u mapi
	objAttrMap := make(map[string][]domain.Attribute)

	// za svaki policy proveri da li trenutno daje dozvolu subjektu
	for _, policy := range resp.Policies {
		objAttrs, ok := objAttrMap[policy.Object.Name()]
		if !ok {
			objAttrs, err = h.getAttributes(policy.Object)
			if err != nil {
				log.Println(err)
				continue
			}
			objAttrMap[policy.Object.Name()] = objAttrs
		}

		hierarchyResp := h.repo.GetPermissionHierarchy(domain.GetPermissionHierarchyReq{
			Subject:        req.Subject,
			Object:         policy.Object,
			PermissionName: policy.PermissionName,
		})
		if hierarchyResp.Error != nil {
			log.Println(hierarchyResp.Error)
			continue
		}

		evalReq := domain.PermissionEvalRequest{
			Subject: subAttrs,
			Object:  objAttrs,
			Env:     req.Env,
		}
		evalResp := hierarchyResp.Hierarchy.Eval(evalReq)
		if authorized(evalResp) {
			granted = append(granted, domain.GrantedPermission{
				PermissionName: policy.PermissionName,
				Object:         policy.Object,
			})
		}
	}

	return domain.GetGrantedPermissionsResp{
		Permissions: granted,
		Error:       nil,
	}
}

func (h EvaluationService) getAttributes(resource domain.Resource) ([]domain.Attribute, error) {
	res := h.repo.GetResource(domain.GetResourceReq{Resource: resource})
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Resource.Attributes, nil
}

func authorized(result domain.EvalResult) bool {
	return result == domain.EvalResultAllowed
}
