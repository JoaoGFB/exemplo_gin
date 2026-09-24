## sugestão de testes

Para testar a API utilizando a extensão REST Client no VSCode ou importar no Insomnia/Postman, pode-se usar os blocos de requisições abaixo, cole em um arquivo chamado testes.http e execute as chamadas na ordem:

```http
### 1- criar sala 1 (capacidade 3)
POST http://localhost:8080/api/v1/salas
Content-Type: application/json

{
    "id": "SALA-01",
    "nome": "Laboratório",
    "capacidade": 3,
    "recursos": ["Computadores", "Projetor", "Ar-condicionado"]
}

### 2- criar sala 2 (capacidade 7)
POST http://localhost:8080/api/v1/salas
Content-Type: application/json

{
    "id": "SALA-02",
    "nome": "Auditório",
    "capacidade": 7,
    "recursos": ["Computadores", "Projetor", "Ar-condicionado"]
}

### 3- criar alunos
POST http://localhost:8080/api/v1/alunos
Content-Type: application/json

{
    "matricula": "001",
    "nome": "Mateus",
    "email": "mateus@email.com"
}

###
POST http://localhost:8080/api/v1/alunos
Content-Type: application/json

{
    "matricula": "002",
    "nome": "Marcos",
    "email": "marcos@email.com"
}

###
POST http://localhost:8080/api/v1/alunos
Content-Type: application/json

{
    "matricula": "003",
    "nome": "Lucas",
    "email": "lucas@email.com"
}

###
POST http://localhost:8080/api/v1/alunos
Content-Type: application/json

{
    "matricula": "004",
    "nome": "João",
    "email": "joaog@email.com"
}

### 4- criar turma
POST http://localhost:8080/api/v1/turmas
Content-Type: application/json

{
    "id": "TURMA-GO",
    "nome": "Backend com GO",
    "disciplina": "Programação Backend",
    "professor": "Prof. Tiago Ravache"
}

### 5- matricular os 4 alunos na turma de Go
POST http://localhost:8080/api/v1/turmas/TURMA-GO/matriculas
Content-Type: application/json

{
    "aluno_id": "001"
}

###
POST http://localhost:8080/api/v1/turmas/TURMA-GO/matriculas
Content-Type: application/json

{
    "aluno_id": "002"
}

###
POST http://localhost:8080/api/v1/turmas/TURMA-GO/matriculas
Content-Type: application/json

{
    "aluno_id": "003"
}

###
POST http://localhost:8080/api/v1/turmas/TURMA-GO/matriculas
Content-Type: application/json

{
    "aluno_id": "004"
}

### 6- testar regra de capacidade (ERRO 422 - a sala 1 só tem 3 lugares, e a turma de Go tem 4)
POST http://localhost:8080/api/v1/turmas/TURMA-GO/alocar
Content-Type: application/json

{
    "sala_id": "SALA-01",
    "dia_semana": "SEGUNDA",
    "horario_inicio": "19:00",
    "horario_termino": "22:00"
}

### 7- alocar turma (sala 2 com 7 lugares)
POST http://localhost:8080/api/v1/turmas/TURMA-GO/alocar
Content-Type: application/json

{
    "sala_id": "SALA-02",
    "dia_semana": "SEGUNDA",
    "horario_inicio": "19:00",
    "horario_termino": "22:00"
}

### 8- verificar as turmas e o status da alocação
GET http://localhost:8080/api/v1/turmas

### 9- verificar a grade da sala alocada
GET http://localhost:8080/api/v1/salas/SALA-02/grade

### 10- consultar dados de um aluno
GET http://localhost:8080/api/v1/alunos/004
