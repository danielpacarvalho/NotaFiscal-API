package controllers

import (
	"database/sql"
	"math/rand"
	"time"

	_ "github.com/nakagami/firebirdsql"
	"github.com/revel/revel"
)

const quixere = "connection string to firebird database"

func EsseRPSExiste(nrrps string, prestador string) bool {
	retorno := false
	var n int
	conn, err := sql.Open("firebirdsql", quixere)
	if err != nil {
		revel.AppLog.Error("Erro durante a conexão com o banco de dados", err)
	}
	defer conn.Close()
	querySelect := "SELECT COUNT(NRRPS)	FROM TB_NFSE WHERE PRESTADOR = %d AND NRRPS = %d '"
	conn.QueryRow(querySelect, prestador, nrrps).Scan(&n)
	//revel.AppLog.Info(querySelect, "Resultado", n)
	if n != 0 {
		retorno = true
	}
	return retorno
}

func CarregarID() int {
	querySelect := "SELECT FIRST 1 ID FROM TB_NFSE ORDER BY ID DESC"
	conn, err := sql.Open("firebirdsql", quixere)
	if err != nil {
		revel.AppLog.Error("Erro durante a conexão com o banco de dados", err)
	}
	defer conn.Close()
	var n int
	conn.QueryRow(querySelect).Scan(&n)
	n++
	return n
}

func CarregarSequencia(prestador string) int {
	querySelect := "SELECT  FIRST 1 SEQUENCIA FROM TB_NFSE WHERE PRESTADOR = %d ORDER BY SEQUENCIA DESC"
	conn, err := sql.Open("firebirdsql", quixere)
	if err != nil {
		revel.AppLog.Error("Erro durante a conexão com o banco de dados", err)
	}
	defer conn.Close()
	var n int
	conn.QueryRow(querySelect, prestador).Scan(&n)
	n++
	return n
}

func GerarVerificador() int {
	rand.Seed(time.Now().UnixNano())
	min := 11111111
	max := 99999999
	return rand.Intn(max-min+1) + min
}

func CarregarAliquotaICSM(prestador string) float32 {
	querySelect := "SELECT ALIQUOTA_ISS	FROM TB_CONTRIBUINTES WHERE COD_CONTRIBUINTE = %d'"
	aliquotaISS := float32(0.0)
	conn, err := sql.Open("firebirdsql", quixere)
	if err != nil {
		revel.AppLog.Error("Erro durante a conexão com o banco de dados", err)
	}
	defer conn.Close()
	conn.QueryRow(querySelect, prestador).Scan(&aliquotaISS)
	return aliquotaISS
}

func GravarNota(id int, sequencia int, dataEmissao string, verificador int, aliquotaISS float64, valorISS float64, jsonData map[string]string) {
	sqlInsert := `INSERT INTO TB_NFSE (ID,
    		                            PRESTADOR,
	        	                        SEQUENCIA,
                		                DATA_EMISSAO, 
                        		        ANO, 
                                		MES, 
                                        NRRPS,
                                        NRNOTA_SUB, 
                                        VALOR_NOTA,
                                        NOME_TOMADOR, 
                                        TCNPJ_CPF, 
                                        TINSCRICAO, 
                                        TLOGRADOURO, 
                                        TNUMERO, 
                                        TCEP, 
                                        TBAIRRO, 
                                        TCIDADE, 
                                        TUF, 
                                        TEMAIL, 
                                        TSUBSTITUTO, 
                                        TIPO_SERVICO, 
                                        LOCAL, 
                                        MODALIDADE, 
                                        VALOR_INSS, 
                                        VALOR_IRRF, 
                                        VALOR_PIS, 
                                        VALOR_COFINS, 
                                        VALOR_CSLL, 
                                        RETIDO_FONTE, 
                                        OBS,
                                        WS)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32)`
	conn, err := sql.Open("firebirdsql", quixere)
	if err != nil {
		revel.AppLog.Error("Erro durante a conexão com o banco de dados", err)
	}
	defer conn.Close()
	conn.QueryRow(sqlInsert, id, jsonData["Prestador"], sequencia, dataEmissao, jsonData["Ano"], jsonData["Mes"], jsonData["Nrrps"],
		jsonData["NrnotaSub"], jsonData["ValorNota"], jsonData["NomeTomador"], jsonData["CnpjCpf"], jsonData["Inscricao"],
		jsonData["Logradouro"], jsonData["Numero"], jsonData["Cep"], jsonData["Bairro"], jsonData["Cidade"], jsonData["UF"],
		jsonData["Email"], " ", jsonData["TipoServico"], jsonData["Local"], jsonData["Modalidade"], jsonData["ValorInss"],
		jsonData["ValorIrrf"], jsonData["ValorPis"], jsonData["ValorCofins"], jsonData["ValorCsll"], jsonData["RetidoNaFonte"],
		jsonData["Obs"], "S")
}
