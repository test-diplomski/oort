package servers

import (
	"log"

	"github.com/c12s/oort/internal/domain"
	"github.com/c12s/oort/internal/mappers/proto"
	"github.com/c12s/oort/internal/services"
	"github.com/c12s/oort/pkg/api"
	"github.com/c12s/oort/pkg/messaging"
)

type AdministratorAsyncServer struct {
	service    services.AdministrationService
	publisher  messaging.Publisher
	subscriber messaging.Subscriber
}

func NewAdministratorAsyncServer(subscriber messaging.Subscriber, publisher messaging.Publisher, service services.AdministrationService) (*AdministratorAsyncServer, error) {
	return &AdministratorAsyncServer{
		service:    service,
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}

func (s *AdministratorAsyncServer) Serve() error {
	return s.subscriber.Subscribe(s.serve)
}

func (s *AdministratorAsyncServer) serve(adminReqMarshalled []byte, replySubject string) {
	adminReq := &api.AdministrationAsyncReq{}
	err := adminReq.Unmarshal(adminReqMarshalled)
	if err != nil {
		log.Println(err)
		return
	}
	var domainResp domain.AdministrationResp
	switch adminReq.Kind {
	case api.AdministrationAsyncReq_CreateResource:
		req := &api.CreateResourceReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.CreateResourceReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.CreateResource(*reqDomain)
	case api.AdministrationAsyncReq_DeleteResource:
		req := &api.DeleteResourceReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.DeleteResourceReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.DeleteResource(*reqDomain)
	case api.AdministrationAsyncReq_PutAttribute:
		req := &api.PutAttributeReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.PutAttributeReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.PutAttribute(*reqDomain)
	case api.AdministrationAsyncReq_DeleteAttribute:
		req := &api.DeleteAttributeReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.DeleteAttributeReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.DeleteAttribute(*reqDomain)
	case api.AdministrationAsyncReq_CreateInheritanceRel:
		req := &api.CreateInheritanceRelReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.CreateInheritanceRelReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.CreateInheritanceRel(*reqDomain)
	case api.AdministrationAsyncReq_DeleteInheritanceRel:
		req := &api.DeleteInheritanceRelReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.DeleteInheritanceRelReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.DeleteInheritanceRel(*reqDomain)
	case api.AdministrationAsyncReq_CreatePolicy:
		req := &api.CreatePolicyReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.CreatePolicyReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.CreatePolicy(*reqDomain)
	case api.AdministrationAsyncReq_DeletePolicy:
		req := &api.DeletePolicyReq{}
		err := req.Unmarshal(adminReq.ReqMarshalled)
		if err != nil {
			log.Println(err)
			return
		}
		reqDomain, err := proto.DeletePolicyReqToDomain(req)
		if err != nil {
			log.Println(err)
			return
		}
		domainResp = s.service.DeletePolicy(*reqDomain)
	default:
		log.Println("unknown request kind")
		return
	}
	resp, err := proto.AdministrationAsyncRespFromDomain(domainResp)
	if err != nil {
		log.Println(err)
		return
	}
	respMarshalled, err := resp.Marshal()
	if err != nil {
		log.Println(err)
		return
	}
	err = s.publisher.Publish(respMarshalled, replySubject)
	if err != nil {
		log.Println(err)
	}
}

func (s *AdministratorAsyncServer) GracefulStop() {
	err := s.subscriber.Unsubscribe()
	if err != nil {
		log.Println(err)
	}
}
