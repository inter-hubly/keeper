//go:generate mockgen -source=template.go -destination=mocks/template_mock.go -package=mocks

package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/keeper/internal/app/domain/kdto"
	"github.com/inter-hubly/keeper/internal/app/gateway"
	"github.com/inter-hubly/keeper/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Template interface {
	SearchTemplates(ctx context.Context, user *hctx.Logged) ([]domain.Template, error)
	SincronizeWhatsAppTemplate(ctx context.Context, user *hctx.Logged) error
	SaveVariable(ctx context.Context, variables []kdto.Variable, templateId string) error
	CountVariables(ctx context.Context, templateId string) (int, error)
}

var (
	_templateServiceOnce sync.Once
	_templateService     *templateService
)

type templateService struct {
	templateRepository repository.Template
	variableRepository repository.Variable
	whatsAppGateway    gateway.WhatsApp
}

func NewTemplate(ctx context.Context) *templateService {
	_templateServiceOnce.Do(func() {
		_templateService = &templateService{
			templateRepository: repository.NewTemplate(ctx),
			variableRepository: repository.NewVariables(ctx),
			whatsAppGateway:    gateway.NewWhatsApp(ctx),
		}
	})
	return _templateService
}

func (s *templateService) SincronizeWhatsAppTemplate(ctx context.Context, user *hctx.Logged) error {
	hlog.Debug(ctx, "templateService.SincronizeWhatsAppTemplate", fmt.Sprintf("Sincronize %s", user.UserId))
	allTemplates, err := s.whatsAppGateway.FindAllTemplate(ctx)
	if err != nil {
		hlog.Error(ctx, "templateService.SincronizeWhatsAppTemplate", err.Error())
		return err
	}
	toSave := make([]domain.Template, 0, len(allTemplates))

	for i := range allTemplates {
		if allTemplates[i].Category != "NAMED" {

			template := domain.Template{
				Name:            strings.ReplaceAll(allTemplates[i].Name, "_", " "),
				Components:      normalizeComponent(allTemplates[i].Components),
				Language:        allTemplates[i].Language,
				ParameterFormat: allTemplates[i].ParameterFormat,
				Status:          domain.TemplateStatus(allTemplates[i].Status),
				Category:        allTemplates[i].Category,
				Slug:            allTemplates[i].Name,
			}
			template.Entity = base.NewBaseEntity(ctx, user)
			toSave = append(toSave, template)
		}
	}

	if err = s.templateRepository.SaveManyTemplate(ctx, toSave); err != nil {
		hlog.Error(ctx, "templateService.SincronizeWhatsAppTemplate", err.Error())
		return err
	}
	return nil
}
func normalizeComponent(components []kdto.ComponentDto) []domain.Component {
	result := make([]domain.Component, 0, len(components))
	for i := range components {
		example := make(map[string][][]string)

		if components[i].Example != nil {
			for k, v := range components[i].Example {
				listOfLists := make([][]string, 0)
				if list, ok := v.([]interface{}); ok {
					for _, inner := range list {
						if innerList, ok := inner.([]interface{}); ok {
							strList := make([]string, 0, len(innerList))
							for _, str := range innerList {
								if s, ok := str.(string); ok {
									strList = append(strList, s)
								}
							}
							listOfLists = append(listOfLists, strList)
						}
					}
				}
				example[k] = listOfLists
			}
		}

		comp := domain.Component{
			Type:    domain.TemplateType(components[i].Type),
			Format:  components[i].Format,
			Text:    components[i].Text,
			Example: example,
		}
		result = append(result, comp)
	}
	return result
}

func (s *templateService) SaveVariable(ctx context.Context, variables []kdto.Variable, templateId string) error {
	hlog.Debug(ctx, "templateService.SaveVariable", fmt.Sprintf("templateId: %s", templateId))
	template, err := s.templateRepository.GetTemplateById(ctx, templateId)
	if err != nil {
		hlog.Error(ctx, "templateService.SaveVariable", err.Error())
		return err
	}

	if err = s.variableRepository.VerifyUserVariables(ctx, variables); err != nil {
		hlog.Error(ctx, "templateService.SaveVariable", err.Error())
		return err
	}
	variablesString := make([]string, 0, len(variables))
	for i := range variables {
		variablesString = append(variablesString, variables[i].Slug)
	}
	template.Variables = variablesString

	if _, err = s.templateRepository.SaveTemplate(ctx, template); err != nil {
		hlog.Error(ctx, "templateService.SaveVariable", err.Error())
		return err
	}
	return nil
}

func (s *templateService) SearchTemplates(ctx context.Context, user *hctx.Logged) ([]domain.Template, error) {
	hlog.Debug(ctx, "templateService.FindAll", fmt.Sprintf("Find All Templates %s", user.UserId))
	return s.templateRepository.SearchTemplates(ctx, user)
}

func (s *templateService) CountVariables(ctx context.Context, templateId string) (int, error) {
	hlog.Debug(ctx, "templateService.CountVariables", fmt.Sprintf("templateId: %s", templateId))
	template, err := s.templateRepository.GetTemplateById(ctx, templateId)
	if err != nil {
		hlog.Error(ctx, "templateService.CountVariables", err.Error())
		return 0, err
	}
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	components := template.Components
	count := 0

	for i := range components {
		matches := re.FindAllStringSubmatch(*components[i].Text, -1)
		if matches != nil {
			count += len(matches)
		}
	}
	return count, nil
}
