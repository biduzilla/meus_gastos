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

### Variáveis de ambiente principais
| Variável                | Descrição                          | Exemplo                          |
|-------------------------|-------------------------------------|-----------------------------------|
| `DB_DSN`                | String de conexão PostgreSQL        | `postgres://user:pass@localhost:5432/meus_gastos?sslmode=disable` |
| `SECURITY_SECRET_KEY`   | Chave secreta para assinar tokens   | `minha_chave_secreta`             |
| `MAIL_HOST` / `MAIL_PORT` / etc. | SMTP para envio de e-mails |                                   |

---

## ▶️ Rodando o projeto

```bash
go mod tidy
go run ./cmd/api
