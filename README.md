## Execução e Teste

Para rodar o server:
```bash 
    go run ./cmd/api
```

Para testar endpoints:
```bash 
    ./test_orders.sh
```

## Arquitetura:
```bash
    order-computations/
├── cmd/api                 → bootstrap
├── internal/domain         → regras de negócio puras
├── internal/application    → orquestração (use cases)
├── internal/interfaces     → delivery (HTTP)
├── internal/infrastructure → detalhes externos (engine, repo)
├── rules/                  → configuração (JSON)
```