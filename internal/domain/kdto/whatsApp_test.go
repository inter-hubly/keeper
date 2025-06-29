package kdto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var resp struct {
		Data []WhatsAppMessageTemplateResponse `json:"data"`
	}

	err := json.Unmarshal([]byte(test), &resp)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp)
}

var test = `
{
  "data" : [ {
    "name" : "realizar_pagamento",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Bom dia!"
    }, {
      "type" : "BODY",
      "text" : "Olá! Tudo bem? Só passando para lembrar sobre o valor de {{1}} que combinamos. Você poderia me dar um retorno sobre o pagamento, por gentileza?",
      "example" : {
        "body_text" : [ [ "100.20" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "Grato!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "1015114033920535"
  }, {
    "name" : "realizar_pagamento_1",
    "previous_category" : "UTILITY",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Seja bem vindo!"
    }, {
      "type" : "BODY",
      "text" : "Aqui você escreve a mensagem com variaveis {{1}}",
      "example" : {
        "body_text" : [ [ "10.00" ] ]
      }
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "MARKETING",
    "correct_category" : "MARKETING",
    "id" : "704424738913490"
  }, {
    "name" : "teste_final",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Seja bem vindo ao teste!"
    }, {
      "type" : "BODY",
      "text" : "Quero informar o valor de sua conta, que é de: {{value}}"
    }, {
      "type" : "FOOTER",
      "text" : "Desde já, grato!"
    } ],
    "language" : "pt_BR",
    "status" : "REJECTED",
    "category" : "UTILITY",
    "id" : "1402734687528571"
  }, {
    "name" : "saimon_tempalte_1",
    "previous_category" : "UTILITY",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Seja bem vindo"
    }, {
      "type" : "BODY",
      "text" : "O valor da sua conta é de {{1}}",
      "example" : {
        "body_text" : [ [ "10.20" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "Desde já grato!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "MARKETING",
    "correct_category" : "MARKETING",
    "id" : "1010064227928407"
  }, {
    "name" : "template_novo_positional_teste",
    "previous_category" : "UTILITY",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Seja bem vindo!"
    }, {
      "type" : "BODY",
      "text" : "Gostaria de informar que voce {{1}} está na nossa lista",
      "example" : {
        "body_text" : [ [ "Saimon" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "Desde já grato!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "MARKETING",
    "correct_category" : "MARKETING",
    "id" : "1361436251671845"
  }, {
    "name" : "template_teste",
    "previous_category" : "UTILITY",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Seja bem vindo!"
    }, {
      "type" : "BODY",
      "text" : "Gostaria de informar que voce {{1}} está na nossa lista",
      "example" : {
        "body_text" : [ [ "Saimon" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "Desde já grato!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "MARKETING",
    "correct_category" : "MARKETING",
    "id" : "1005043494906370"
  }, {
    "name" : "primeiro_teste",
    "previous_category" : "UTILITY",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Boas-vindas"
    }, {
      "type" : "BODY",
      "text" : "Olá! Seja muito bem-vindo(a)!"
    }, {
      "type" : "FOOTER",
      "text" : "É um prazer ter você por aqui."
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "MARKETING",
    "correct_category" : "MARKETING",
    "id" : "890606413147363"
  }, {
    "name" : "teste_tenant",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Bom dia!"
    }, {
      "type" : "BODY",
      "text" : "Gostaria de informar que seu processo se encaminha pelas vias de fato"
    }, {
      "type" : "FOOTER",
      "text" : "Desde já, grato!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "558729847333553"
  }, {
    "name" : "rosseto_advogados",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Bom dia!"
    }, {
      "type" : "BODY",
      "text" : "Gostaria de informar que seu processo está em execucao"
    }, {
      "type" : "FOOTER",
      "text" : "Muito obrigado!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "679081781222440"
  }, {
    "name" : "cobranca_mensal_1",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Lembrete de Pagamento"
    }, {
      "type" : "BODY",
      "text" : "Oi, {{1}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês de R$ {{2}} 📝. Assim a gente mantém tudo certinho e posso continuar te ajudando a arrasar no inglês! 🚀 Qualquer dúvida, é só me chamar! 😄",
      "example" : {
        "body_text" : [ [ "Saimon", "20,20" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "See ya!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "942970764616888"
  }, {
    "name" : "teste_inicializacao",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Bom dia!"
    }, {
      "type" : "BODY",
      "text" : "Esse mês sua conta venceu dia 2 de abril"
    }, {
      "type" : "FOOTER",
      "text" : "desde já! grato!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "1183668785983556"
  }, {
    "name" : "grande_teste",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "teste"
    }, {
      "type" : "BODY",
      "text" : "teste"
    }, {
      "type" : "FOOTER",
      "text" : "teste"
    } ],
    "language" : "pt_BR",
    "status" : "REJECTED",
    "category" : "UTILITY",
    "id" : "664039912942399"
  }, {
    "name" : "cobranca_mensal",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Lembrete de Pagamento"
    }, {
      "type" : "BODY",
      "text" : "Oi, {{1}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês de R$ {{2}} 📝. Assim a gente mantém tudo certinho e posso continuar te ajudando a arrasar no inglês! 🚀 Qualquer dúvida, é só me chamar! 😄",
      "example" : {
        "body_text" : [ [ "Saimon", "20,20" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "See ya!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "sub_category" : "CUSTOM",
    "id" : "908341171131680"
  }, {
    "name" : "inicial_1",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Hello!"
    }, {
      "type" : "BODY",
      "text" : "Oi, {{1}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês de R$ {{2}} 📝. Assim a gente mantém tudo certinho e posso continuar te ajudando a arrasar no inglês! 🚀 Qualquer dúvida, é só me chamar! 😄",
      "example" : {
        "body_text" : [ [ "Joao", "123,12" ] ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "See ya!"
    } ],
    "language" : "en_US",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "library_template_name" : "delivery_confirmation_1",
    "id" : "977234657580494"
  }, {
    "name" : "cobranca",
    "parameter_format" : "NAMED",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Lembrete de Pagamento"
    }, {
      "type" : "BODY",
      "text" : "Oi, {{owner_name}}! Tudo bem? 😊 Só passando para lembrar sobre o pagamento da aula de inglês de R$ {{owner_value}} 📝. Assim a gente mantém tudo certinho e posso continuar te ajudando a arrasar no inglês! 🚀 Qualquer dúvida, é só me chamar! 😄",
      "example" : {
        "body_text_named_params" : [ {
          "param_name" : "owner_name",
          "example" : "owner_name"
        }, {
          "param_name" : "owner_value",
          "example" : "owner_value"
        } ]
      }
    }, {
      "type" : "FOOTER",
      "text" : "See Ya!"
    } ],
    "language" : "pt_BR",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "535167829552519"
  }, {
    "name" : "inicial",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "BODY",
      "text" : "Bom Dia!"
    } ],
    "language" : "pt_BR",
    "status" : "REJECTED",
    "category" : "UTILITY",
    "sub_category" : "CUSTOM",
    "id" : "1529019937815243"
  }, {
    "name" : "hello_world",
    "parameter_format" : "POSITIONAL",
    "components" : [ {
      "type" : "HEADER",
      "format" : "TEXT",
      "text" : "Hello World"
    }, {
      "type" : "BODY",
      "text" : "Welcome and congratulations!! This message demonstrates your ability to send a WhatsApp message notification from the Cloud API, hosted by Meta. Thank you for taking the time to test with us."
    }, {
      "type" : "FOOTER",
      "text" : "WhatsApp Business Platform sample message"
    } ],
    "language" : "en_US",
    "status" : "APPROVED",
    "category" : "UTILITY",
    "id" : "2799264690254265"
  } ],
  "paging" : {
    "cursors" : {
      "before" : "MAZDZD",
      "after" : "MjQZD"
    }
  }
}`
