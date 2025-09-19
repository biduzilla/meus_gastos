# 💰 Meus Gastos – API

API REST em Go para gerenciamento de **usuários**, **categorias** e **transações financeiras**.  
Inclui autenticação via **JWT**, controle de taxa (**rate limiting**), CORS configurável e métricas de desempenho.

---

## 🚀 Funcionalidades
- Registro e ativação de usuários com código de confirmação.
- Autenticação via token JWT.
- CRUD de **categorias** (ex.: Alimentação, Lazer).
- CRUD de **transações** vinculadas a categorias.
- Filtros de data, descrição e tipo de categoria.
- Métricas expostas em `/debug/vars`.

---

## 🛠️ Requisitos
- Go 1.21+
- PostgreSQL 13+
- Variáveis de ambiente definidas (ou arquivo de configuração compatível).

---

## ⚙️ Configuração

### Variáveis de Ambiente

Defina as variáveis abaixo (por exemplo em um arquivo `.env` ou no sistema):

| Variável                | Descrição                                      | Exemplo Seguro                                           |
|-------------------------|------------------------------------------------|----------------------------------------------------------|
| `SERVER_PORT`           | Porta em que a API vai rodar                    | `4000`                                                  |
| `SERVER_TIMEOUT_READ`   | Timeout de leitura do servidor                  | `3s`                                                    |
| `SERVER_TIMEOUT_WRITE`  | Timeout de escrita do servidor                  | `5s`                                                    |
| `SERVER_TIMEOUT_IDLE`   | Timeout de conexão ociosa                       | `5s`                                                    |
| `SERVER_DEBUG`          | Ativa modo debug                                | `true`                                                  |
| `POSTGRES_USER`         | Usuário do banco PostgreSQL                     | `meusgastos_user`                                       |
| `POSTGRES_PASSWORD`     | Senha do banco PostgreSQL                       | `sua_senha_segura`                                      |
| `POSTGRES_DB`           | Nome do banco de dados                          | `meusgastos_db`                                         |
| `POSTGRES_PORT`         | Porta do PostgreSQL                              | `5432`                                                  |
| `DB_DSN`                | String de conexão (usa variáveis acima)          | `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable` |
| `DB_MAX_OPEN_CONNS`     | Máx. conexões abertas                             | `25`                                                   |
| `DB_MAX_IDLE_CONNS`     | Máx. conexões inativas                            | `25`                                                   |
| `DB_MAX_IDLE_TIME`      | Tempo máx. que uma conexão pode ficar inativa     | `15m`                                                  |
| `LIMITER_RPS`           | Requisições por segundo permitidas               | `2`                                                    |
| `LIMITER_BURST`         | Burst máximo de requisições                       | `4`                                                    |
| `LIMITER_ENABLED`       | Ativa/desativa rate limiting                      | `true`                                                 |
| `EMAIL_HOST`            | Host do servidor SMTP                             | `smtp.seuprovedor.com`                                  |
| `EMAIL_POST`            | Porta SMTP                                        | `587`                                                  |
| `EMAIL_USERNAME`        | Usuário SMTP                                      | `usuario_smtp`                                         |
| `EMAIL_PASSWORD`        | Senha SMTP                                        | `senha_smtp_segura`                                    |
| `SECRET_KEY`            | Chave secreta para assinar tokens (JWT)           | `uma_chave_secreta_bem_grande_e_aleatoria`             |


---

## ▶️ Rodando o projeto

```bash
go mod tidy
go run ./cmd/api
