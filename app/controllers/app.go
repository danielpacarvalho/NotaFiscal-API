package controllers

import (
	"strconv"
	"time"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	mensagem := "Serviço de Nota Fiscal Eletrônica, Prefeitura de Quixerê"
	return c.Render(mensagem)
	//return "Sistema de Nota Fiscal Eletrônica, Prefeitura de Quixerê"
}

func (c App) Gravar() revel.Result {
	//Esse eh o arquivo recebido
	var jsonData map[string]string
	c.Params.BindJSON(&jsonData)

	//Checar se está tudo ok
	mensagem := " "
	check := true
	if len(jsonData["NomeTomador"]) == 0 {
		mensagem += "O nome do tomador não pode ser vazio "
		check = false
	}

	if len(jsonData["ValorNota"]) == 0 {
		mensagem += "O valor da nota não pode ser vazio "
		check = false
	}

	if len(jsonData["CnpjCpf"]) == 0 {
		mensagem += "O número de CPF ou CNPJ deve ser informado "
		check = false
	}

	if len(jsonData["Logradouro"]) == 0 {
		mensagem += "O Logradouro não pode ser vazio "
		check = false
	}

	if len(jsonData["Numero"]) == 0 {
		mensagem += "O número no endereço não pode ser vazio "
		check = false
	}

	if len(jsonData["Bairro"]) == 0 {
		mensagem += "O nome do bairro não pode ser vazio "
		check = false
	}

	if len(jsonData["Cidade"]) == 0 {
		mensagem += "O nome do município não pode ser vazio "
		check = false
	}

	if len(jsonData["UF"]) == 0 {
		mensagem += "O nome do Estado não pode ser vazio "
		check = false
	}
	if len(jsonData["Local"]) == 0 {
		mensagem += "O local da prestação de serviço não pode ser vazio "
		check = false
	}

	if len(jsonData["RetidoNaFonte"]) == 0 {
		mensagem += "O campo RetidoNaFonte não pode ser vazio "
		check = false
	}
	modalidade, _ := strconv.Atoi(jsonData["Modalidade"])
	if modalidade <= 0 || modalidade > 4 {
		mensagem += "A modalidade deve ser informada, verifique se o numero da modalidade está correto "
		check = false
	}

	if len(jsonData["RetidoNaFonte"]) == 0 {
		mensagem += "O campo RetidoNaFonte não pode ser vazio "
		check = false
	}

	if len(jsonData["Prestador"]) == 0 {
		mensagem += "O campo de código do Prestador não pode ser vazio "
		check = false
	} else {
		if len(jsonData["Nrrps"]) == 0 {
			mensagem += "O número de RPS não pode ser vazio "
			check = false
		} else {
			if EsseRPSExiste(jsonData["Nrrps"], jsonData["Prestador"]) {
				mensagem += "O número de RPS informado já foi cadastrado "
				check = false
			}
		}
	}

	//Se tudo ok, partir para gravar, se não retornar dizendo o que está errado
	if check {
		//Carregar dados do prestador
		id := CarregarID()
		sequencia := CarregarSequencia(jsonData["Prestador"])
		dataEmissao := time.Now()
		verificador := GerarVerificador()

		//Carrega a aliquotaISS do prestador e tranformar em %
		aliquotaISS := float64(CarregarAliquotaICSM(jsonData["Prestador"]) / 100)
		valorNota, _ := strconv.ParseFloat(jsonData["ValorNota"], 32)
		descontoPrevistaEmLei, _ := strconv.ParseFloat(jsonData["DescontoPrevistaEmLei"], 32)
		valorDescincond, _ := strconv.ParseFloat(jsonData["ValorDescincond"], 32)
		valorISS := (valorNota - descontoPrevistaEmLei - valorDescincond) * aliquotaISS

		//Gravar a nota
		GravarNota(id, sequencia, dataEmissao.String(), verificador, aliquotaISS, valorISS, jsonData)

		//Enviar mensagem de retorno
		mensagem = "Gravação de NF Nº " + strconv.Itoa(sequencia) + " executada com sucesso"
		return c.RenderText(mensagem)

	} else {
		//Não deu certo vamos retornar com o que está errado
		return c.RenderText(mensagem)
	}
}
