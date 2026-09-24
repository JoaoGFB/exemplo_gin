package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Sala struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   []string `json:"recursos"`
}
type Aluno struct {
	Matricula string `json:"matricula"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
}
type Alocacao struct {
	SalaID  string `json:"sala_id"`
	Dia     string `json:"dia_semana"`
	Inicio  string `json:"horario_inicio"`
	Termino string `json:"horario_termino"`
}
type Turma struct {
	ID         string    `json:"id"`
	Nome       string    `json:"nome"`
	Disciplina string    `json:"disciplina"`
	Professor  string    `json:"professor"`
	Alunos     []string  `json:"alunos_matriculados"`
	Alocacao   *Alocacao `json:"alocacao,omitempty"`
}

var (
	dbSalas  = make(map[string]Sala)
	dbAlunos = make(map[string]Aluno)
	dbTurmas = make(map[string]Turma)
)

func temSobreposicao(inicio1, termino1, inicio2, termino2 string) bool {
	return inicio1 < termino2 && termino1 > inicio2
}
func alunoTemConflitoHorario(matricula, dia, inicio, termino string) bool {
	for _, turma := range dbTurmas {
		if turma.Alocacao != nil && turma.Alocacao.Dia == dia {
			for _, mat := range turma.Alunos {
				if mat == matricula {
					if temSobreposicao(turma.Alocacao.Inicio, turma.Alocacao.Termino, inicio, termino) {
						return true
					}
				}
			}
		}
	}
	return false
}
func main() {

	r := gin.New()

	// Uso dos Middlewares globais nativos e personalizados
	r.Use(gin.Recovery())

	// 4. Mapeamento de Rotas sob Grupo Versionado
	v1 := r.Group("/api/v1")
	{
		// Monitoramento da API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		// Domínio de Turmas (Classes)
		//v1.POST("/turmas", turmaHandler.CriarTurma)
		//v1.GET("/turmas", turmaHandler.ListarTurmas)
		//v1.POST("/turmas/:id/alocar", turmaHandler.AlocarSala)

		v1.POST("/salas", func(c *gin.Context) {
			var sala Sala
			if err := c.ShouldBindJSON(&sala); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos"})
				return
			}
			if sala.Capacidade <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"erro": "capacidade deve ser maior que 0"})
				return
			}
			dbSalas[sala.ID] = sala
			c.JSON(http.StatusCreated, sala)
		})
		v1.GET("/salas", func(c *gin.Context) {
			var lista []Sala
			for _, s := range dbSalas {
				lista = append(lista, s)
			}
			c.JSON(http.StatusOK, lista)
		})
		v1.GET("/salas/:id/grade", func(c *gin.Context) {
			salaID := c.Param("id")
			if _, existe := dbSalas[salaID]; !existe {
				c.JSON(http.StatusNotFound, gin.H{"erro": "sala não encontrada"})
				return
			}
			var gradeUso []gin.H
			for _, t := range dbTurmas {
				if t.Alocacao != nil && t.Alocacao.SalaID == salaID {
					gradeUso = append(gradeUso, gin.H{
						"turma_id":       t.ID,
						"nome_turma":     t.Nome,
						"dia_semana":     t.Alocacao.Dia,
						"horario_inicio": t.Alocacao.Inicio,
						"horario_fim":    t.Alocacao.Termino,
					})
				}
			}
			c.JSON(http.StatusOK, gradeUso)
		})
		v1.POST("/alunos", func(c *gin.Context) {
			var aluno Aluno
			if err := c.ShouldBindJSON(&aluno); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos"})
				return
			}
			dbAlunos[aluno.Matricula] = aluno
			c.JSON(http.StatusCreated, aluno)
		})
		v1.GET("/alunos", func(c *gin.Context) {
			var lista []Aluno
			for _, a := range dbAlunos {
				lista = append(lista, a)
			}
			c.JSON(http.StatusOK, lista)
		})
		v1.GET("/alunos/:id", func(c *gin.Context) {
			matricula := c.Param("id")
			aluno, existe := dbAlunos[matricula]
			if !existe {
				c.JSON(http.StatusNotFound, gin.H{"erro": "aluno não encontrado"})
				return
			}
			c.JSON(http.StatusOK, aluno)
		})
		v1.POST("/turmas", func(c *gin.Context) {
			var turma Turma
			if err := c.ShouldBindJSON(&turma); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos"})
				return
			}
			turma.Alunos = []string{}
			dbTurmas[turma.ID] = turma
			c.JSON(http.StatusCreated, turma)
		})
		v1.GET("/turmas", func(c *gin.Context) {
			var listaFormatada []gin.H
			for _, t := range dbTurmas {
				listaFormatada = append(listaFormatada, gin.H{
					"id":                t.ID,
					"nome":              t.Nome,
					"disciplina":        t.Disciplina,
					"professor":         t.Professor,
					"quantidade_alunos": len(t.Alunos),
					"status_alocacao":   t.Alocacao,
				})
			}
			c.JSON(http.StatusOK, listaFormatada)
		})
		v1.POST("/turmas/:id/matriculas", func(c *gin.Context) {
			turmaID := c.Param("id")
			var req struct {
				AlunoID string `json:"aluno_id"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos"})
				return
			}
			turma, turmaExiste := dbTurmas[turmaID]
			_, alunoExiste := dbAlunos[req.AlunoID]
			if !turmaExiste || !alunoExiste {
				c.JSON(http.StatusNotFound, gin.H{"erro": "turma ou aluno não econtrado"})
				return
			}
			for _, mat := range turma.Alunos {
				if mat == req.AlunoID {
					c.JSON(http.StatusConflict, gin.H{"erro": "aluno já matriculado neste turma"})
					return
				}
			}
			if turma.Alocacao != nil {
				sala := dbSalas[turma.Alocacao.SalaID]
				if len(turma.Alunos) >= sala.Capacidade {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "capacidade máxima da sala atingida"})
					return
				}
				if alunoTemConflitoHorario(req.AlunoID, turma.Alocacao.Dia, turma.Alocacao.Inicio, turma.Alocacao.Termino) {
					c.JSON(http.StatusConflict, gin.H{"erro": "conflito de horário na grade do aluno"})
					return
				}
			}
			turma.Alunos = append(turma.Alunos, req.AlunoID)
			dbTurmas[turmaID] = turma
			c.JSON(http.StatusOK, gin.H{"messagem": "aluno matriculado com sucesso"})
		})
		v1.GET("/turmas/:id/matriculas", func(c *gin.Context) {
			turmaID := c.Param("id")
			turma, existe := dbTurmas[turmaID]
			if !existe {
				c.JSON(http.StatusNotFound, gin.H{"erro": "turma não encontrada"})
				return
			}
			var alunosMatriculados []Aluno
			for _, mat := range turma.Alunos {
				alunosMatriculados = append(alunosMatriculados, dbAlunos[mat])
			}
			c.JSON(http.StatusOK, alunosMatriculados)
		})
		v1.POST("/turmas/:id/alocar", func(c *gin.Context) {
			turmaID := c.Param("id")
			var aloc Alocacao
			if err := c.ShouldBindJSON(&aloc); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos"})
				return
			}
			turma, turmaExiste := dbTurmas[turmaID]
			sala, salaExiste := dbSalas[aloc.SalaID]
			if !turmaExiste || !salaExiste {
				c.JSON(http.StatusNotFound, gin.H{"erro": "turma ou sala não encontrada"})
				return
			}
			if len(turma.Alunos) > sala.Capacidade {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "a capacidade da sala é menor que a quantidade de alunos"})
				return
			}
			for _, t := range dbTurmas {
				if t.ID != turmaID && t.Alocacao != nil && t.Alocacao.SalaID == aloc.SalaID && t.Alocacao.Dia == aloc.Dia {
					if temSobreposicao(t.Alocacao.Inicio, t.Alocacao.Termino, aloc.Inicio, aloc.Termino) {
						c.JSON(http.StatusConflict, gin.H{"erro": "a sala já possui turma alocada neste horário"})
						return
					}
				}
			}
			for _, matricula := range turma.Alunos {
				if alunoTemConflitoHorario(matricula, aloc.Dia, aloc.Inicio, aloc.Termino) {
					c.JSON(http.StatusConflict, gin.H{"erro": "um ou mais alunos possuem conflito de horário com esta alocação"})
					return
				}
			}
			turma.Alocacao = &aloc
			dbTurmas[turmaID] = turma
			c.JSON(http.StatusOK, gin.H{"mensagem": "sala alocada com sucesso"})
		})
	}

	r.Run(":8080")
}
