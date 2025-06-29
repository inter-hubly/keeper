package templates

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/keeper/internal/templates/mocks"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

type allMock struct {
	templateRepository *mocks.MockRepository
}

func TestTemplateService(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	allMocks := allMock{
		templateRepository: mocks.NewMockRepository(ctrl),
	}

	templateServiceTest := templateService{
		templateRepository: allMocks.templateRepository,
	}

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "Need verify all matches",
			auxFunc: func() {
				text := "Need to find one {{1}}"
				allMocks.templateRepository.EXPECT().GetTemplateById(gomock.Any(), gomock.Any()).Return(
					&domain.Template{
						Components: []domain.Component{
							{
								Text: &text,
							},
							{
								Text: &text,
							},
							{
								Text: &text,
							},
						},
					}, nil)

				variables, err := templateServiceTest.CountVariables(ctx, "templateId")
				assert.NoError(t, err)
				assert.Equal(t, 3, variables)
			},
		},
		{
			testName: "Need verify all matches",
			auxFunc: func() {
				text := "Need to find one {{1}} {{2}}"
				allMocks.templateRepository.EXPECT().GetTemplateById(gomock.Any(), gomock.Any()).Return(
					&domain.Template{
						Components: []domain.Component{
							{
								Text: &text,
							},
							{
								Text: &text,
							},
							{
								Text: &text,
							},
						},
					}, nil)

				variables, err := templateServiceTest.CountVariables(ctx, "templateId")
				assert.NoError(t, err)
				assert.Equal(t, 6, variables)
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			v.auxFunc()
		})
	}
}
