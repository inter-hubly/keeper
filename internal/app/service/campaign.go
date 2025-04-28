package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Campaign interface {
	SaveCampaign(ctx context.Context, loggedUser *hctx.Logged, campaignDto *kdto.Campaign) (*entity.Campaign, error)
	GetCampaign(ctx context.Context, user *hctx.Logged) (*entity.Campaign, error)
	StartCampaign(ctx context.Context, user *hctx.Logged, campaignId string) error
	ListCampaign(ctx context.Context, user *hctx.Logged) ([]entity.Campaign, error)
}

type campaignService struct {
	campaignRepository repository.Campaign
	contactRepository  repository.Contact
	variableRepository repository.Variable
	templateRepository repository.Template
	broker             broker.Connection
}

var (
	_campaignServiceOnce sync.Once
	_campaignService     *campaignService
)

func NewCampaign(ctx context.Context) *campaignService {
	_campaignServiceOnce.Do(func() {
		_campaignService = &campaignService{
			campaignRepository: repository.NewCampaign(ctx),
			contactRepository:  repository.NewContact(ctx),
			variableRepository: repository.NewVariables(ctx),
			templateRepository: repository.NewTemplate(ctx),
			broker:             broker.GetConnection(),
		}
	})
	return _campaignService
}

func (c *campaignService) StartCampaign(ctx context.Context, user *hctx.Logged, campaignId string) error {
	hlog.Debug(ctx, "campaignService.StartCampaign", fmt.Sprintf("start campaign %s", user.Username))
	campaignStruct := struct {
		Id string `json:"id"`
	}{
		Id: campaignId,
	}
	campaignBytes, err := json.Marshal(campaignStruct)
	if err != nil {
		hlog.Error(ctx, "campaignService.StartCampaign", fmt.Sprintf("error marshalling campaign struct: %v", err))
		return err
	}

	return c.broker.Publish(ctx, "campaign.init", campaignBytes)
}

func (c *campaignService) GetCampaign(ctx context.Context, user *hctx.Logged) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignService.GetCampaign", fmt.Sprintf("found campaign for %s", user.Username))
	campaign, err := c.campaignRepository.GetCampaignById(ctx, user.Tenant)
	if err != nil {
		hlog.Error(ctx, "campaignService.GetCampaign", fmt.Sprintf("error getting campaign by id: %v", err))
		return nil, err
	}
	return campaign, nil
}

func (c *campaignService) SaveCampaign(ctx context.Context, loggedUser *hctx.Logged, campaignDto *kdto.Campaign) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignService.SaveCampaign", fmt.Sprintf("%s saving campaign with dto %s", loggedUser.Username, campaignDto.Name))
	contacts, err := c.contactRepository.FindContacts(ctx, loggedUser.Tenant)
	if err != nil {
		hlog.Error(ctx, "campaignService.SaveCampaign", fmt.Sprintf("Contact Repository Error: %v", err))
		return nil, err
	}
	if !c.containsAll(contacts, campaignDto.ContactsID) {
		hlog.Error(ctx, "campaignService.SaveCampaign", "this phone does not contain contacts")
		return nil, errors.New("this phone does not contain contacts")
	}
	allVariables, err := c.variableRepository.GetVariables(ctx)
	if err != nil {
		hlog.Error(ctx, "campaignService.SaveCampaign", fmt.Sprintf("VariableRepository Error: %v", err))
		return nil, err
	}
	if !c.containsAllVariables(allVariables, campaignDto.Variables) {
		hlog.Error(ctx, "campaignService.SaveCampaign", "this phone does not contain variables")
		return nil, errors.New("this phone does not contain this variables")
	}

	templateEntity, err := c.templateRepository.GetTemplateById(ctx, campaignDto.TemplateId)
	if err != nil {
		hlog.Error(ctx, "campaignService.SaveCampaign", fmt.Sprintf("Template Repository Error: %v", err))
		return nil, err
	}
	campaignDb := entity.Campaign{
		Name: campaignDto.Name,
		Template: base.TemplateInfo{
			Id:       templateEntity.Id,
			Name:     templateEntity.Slug,
			Language: templateEntity.Language,
			Message:  templateEntity.GetComponentMessages(),
		},
		ContactsId: campaignDto.ContactsID,
		IaContext:  campaignDto.IaContext,
		Variables:  campaignDto.Variables,
		Entity:     base.NewBaseEntity(ctx, loggedUser),
		Flows:      campaignDto.Flows,
	}

	campaign, err := c.campaignRepository.SaveCampaign(ctx, &campaignDb)
	if err != nil {
		hlog.Error(ctx, "campaignService.SaveCampaign", fmt.Sprintf("SaveCampaign Error: %v", err))
		return nil, err
	}
	return campaign, nil
}

func (c *campaignService) containsAll(list1 []domain.Contact, list2 []string) bool {
	elements := make(map[string]bool)

	for _, item := range list2 {
		elements[item] = true
	}

	for _, item := range list1 {
		if _, ok := elements[item.Phone]; ok {
			return false
		}
	}

	return true
}

func (c *campaignService) containsAllVariables(list1 map[string]domain.SingleVariable, list2 []valueobject.Pair[string, string]) bool {
	elements := make(map[string]bool)

	for _, item := range list2 {
		elements[item.Key] = true
	}

	for _, item := range list1 {
		if _, ok := elements[item.Slug]; ok {
			return false
		}
	}

	return true
}

func (c *campaignService) ListCampaign(ctx context.Context, user *hctx.Logged) ([]entity.Campaign, error) {
	hlog.Debug(ctx, "campaignService.ListCampaign", fmt.Sprintf("list campaign %s", user.Username))
	return c.campaignRepository.ListCampaign(ctx)
}
