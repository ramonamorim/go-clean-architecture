# Go Expert - Clean Architecture

## Proposta

Criar um usecase para retornar uma lista das __orders__ criadas e ter a possibilidade de lista-las por REST, gRPC e GraphQL.

## Setup 

Clone o repositório:

```bash
git clone https://github.com/ramonamorim/go-clean-architecture.git
```

O projeto está configurado para usar imagens compatíveis com arquitetura arm64, caso use outra arquitetura será necessário ajustar o arquivo docker-compose.yaml.

Para executar o projeto, é necessário ter o docker e o docker-compose instalados. E execute o comando abaixo para subir as imagens:

```bash
docker-compose up -d
```
ou

```bash
make up
```

Os serviços estarão rodando nos endereços:

- Rest em http://localhost:8000:

  - Use os arquivos na pasta `/api` 
  (https://marketplace.visualstudio.com/items?itemName=anweber.vscode-httpyac)

- GraphQL em http://localhost:8080:

  - Use o template abaixo:

    ```graphql
    mutation createOrder {
      createOrder(input: { id: "id", Price: 100.0, Tax: 5 }) {
        id
        Price
        Tax
        FinalPrice
      }
    }

    query orders {
      listOrders {
        id
        Price
        Tax
        FinalPrice
      }
    }
    ```

- gRPC na porta 50051:
  - Será necessário ter o [evans](https://github.com/ktr0731/evans?tab=readme-ov-file#installation) instalado.
  - Feito isso pode executar os comando a seguir:
    ```bash
    evans -r repl
    service OrderService
    call CreateOrder - criar 1 pedido
    call ListOrders - lista os pedidos criados
    ```
