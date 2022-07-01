# dm-go-ft
Repositório para a biblioteca de feature toggle para projetos GO

# Como usar
A biblioteca atualmente utiliza feature toggles que devem ser definidas diretamente no código, seja hard-coded ou em variáveis de ambiente.
Ela também, por enquanto, apenas suporta feature toggles "true" ou "false".

## instalação
```bash
go get github.com/delivery-much/dm-go-ft
```

## Configuração
A biblioteca deve ser instanciada através do método `Init`.
Este método receberá uma série de pares "chave/valor" contendo os feature toggles.

Example:
```go
import "github.com/delivery-much/dm-go-ft/featuretoggle"

...

// pode ser hard-coded ou uma variável de ambiente
myFeatureToggleVal := true

featuretoggle.Init(map[string]interface{}{
  "MyKey": myFeatureToggleVal,
})
```

> de preferência, o Init deverá ser chamado uma vez, no main do projeto.

## Uso
Após a biblioteca ter sido instanciada, pode-se chamar a biblioteca de qualquer ponto do código.
Pode-se usar então o método `IsEnabled` para verificar o feature toggle instanciado previamente.

Example:
```go
import "github.com/delivery-much/dm-go-ft/featuretoggle"

...

if featuretoggle.IsEnabled("MyKey") {
  // do something
}
```
